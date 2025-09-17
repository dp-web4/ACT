#!/usr/bin/env node

/**
 * Swarm Veracity Agent
 * Responsible for verifying that claimed functions, types, and methods actually exist
 * This is the "truth detector" for the swarm - prevents false claims from propagating
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

class VeracityAgent {
  constructor(config = {}) {
    this.name = config.name || 'Veracity-Agent';
    this.lctId = `lct:veracity:${Date.now()}`;
    this.witnessLog = config.witnessLog || './swarm-memory/witness/veracity.jsonl';
    this.verifications = [];
  }

  /**
   * Verify that a type exists in a Go package
   */
  verifyTypeExists(packagePath, typeName) {
    const verification = {
      claim: `Type ${typeName} exists in ${packagePath}`,
      type: 'type_existence',
      verified: false,
      evidence: {}
    };

    try {
      // Find all .go files in the package
      const files = this.findGoFiles(packagePath);
      
      for (const file of files) {
        const content = fs.readFileSync(file, 'utf8');
        
        // Look for type declaration
        const typeRegex = new RegExp(`type\\s+${typeName}\\s+(struct|interface)`, 'g');
        if (typeRegex.test(content)) {
          verification.verified = true;
          verification.evidence = {
            file: path.relative(packagePath, file),
            found: true
          };
          break;
        }
      }
      
      if (!verification.verified) {
        verification.evidence = {
          searched: files.length,
          found: false
        };
      }
    } catch (error) {
      verification.evidence = { error: error.message };
    }

    this.recordVerification(verification);
    return verification;
  }

  /**
   * Verify that a function exists in a package
   */
  verifyFunctionExists(packagePath, functionName, receiverType = null) {
    const verification = {
      claim: receiverType 
        ? `Method ${functionName} exists on ${receiverType}` 
        : `Function ${functionName} exists in ${packagePath}`,
      type: 'function_existence',
      verified: false,
      evidence: {}
    };

    try {
      const files = this.findGoFiles(packagePath);
      
      for (const file of files) {
        const content = fs.readFileSync(file, 'utf8');
        
        // Look for function/method declaration
        let funcRegex;
        if (receiverType) {
          // Method on a type
          funcRegex = new RegExp(`func\\s+\\([^)]*${receiverType}[^)]*\\)\\s+${functionName}\\s*\\(`, 'g');
        } else {
          // Standalone function
          funcRegex = new RegExp(`func\\s+${functionName}\\s*\\(`, 'g');
        }
        
        if (funcRegex.test(content)) {
          verification.verified = true;
          verification.evidence = {
            file: path.relative(packagePath, file),
            found: true,
            receiver: receiverType || 'none'
          };
          break;
        }
      }
      
      if (!verification.verified) {
        verification.evidence = {
          searched: files.length,
          found: false
        };
      }
    } catch (error) {
      verification.evidence = { error: error.message };
    }

    this.recordVerification(verification);
    return verification;
  }

  /**
   * Verify that a field exists on a struct
   */
  verifyFieldExists(packagePath, structName, fieldName) {
    const verification = {
      claim: `Field ${fieldName} exists on ${structName}`,
      type: 'field_existence',
      verified: false,
      evidence: {}
    };

    try {
      const files = this.findGoFiles(packagePath);
      
      for (const file of files) {
        const content = fs.readFileSync(file, 'utf8');
        
        // Find the struct definition
        const structRegex = new RegExp(`type\\s+${structName}\\s+struct\\s*{([^}]+)}`, 's');
        const structMatch = content.match(structRegex);
        
        if (structMatch) {
          const structBody = structMatch[1];
          // Look for the field
          const fieldRegex = new RegExp(`\\b${fieldName}\\b`, 'g');
          if (fieldRegex.test(structBody)) {
            verification.verified = true;
            verification.evidence = {
              file: path.relative(packagePath, file),
              struct: structName,
              field: fieldName,
              found: true
            };
            break;
          }
        }
      }
      
      if (!verification.verified) {
        verification.evidence = {
          searched: files.length,
          struct: structName,
          field: fieldName,
          found: false
        };
      }
    } catch (error) {
      verification.evidence = { error: error.message };
    }

    this.recordVerification(verification);
    return verification;
  }

  /**
   * Verify that an import is used in a file
   */
  verifyImportUsed(filePath, importPath) {
    const verification = {
      claim: `Import ${importPath} is used in ${filePath}`,
      type: 'import_usage',
      verified: false,
      evidence: {}
    };

    try {
      const content = fs.readFileSync(filePath, 'utf8');
      
      // Check if import exists
      const importRegex = new RegExp(`import.*"${importPath.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}"`, 's');
      const hasImport = importRegex.test(content);
      
      if (hasImport) {
        // Extract package alias if any
        const aliasMatch = content.match(new RegExp(`(\\w+)\\s+"${importPath.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}"`));
        const packageName = aliasMatch ? aliasMatch[1] : path.basename(importPath);
        
        // Check if package is actually used
        const usageRegex = new RegExp(`\\b${packageName}\\.\\w+`, 'g');
        const isUsed = usageRegex.test(content);
        
        verification.verified = isUsed;
        verification.evidence = {
          imported: true,
          used: isUsed,
          packageName
        };
      } else {
        verification.evidence = {
          imported: false,
          used: false
        };
      }
    } catch (error) {
      verification.evidence = { error: error.message };
    }

    this.recordVerification(verification);
    return verification;
  }

  /**
   * Verify a swarm claim before execution
   */
  async verifySwarmClaim(claim) {
    const verification = {
      claim: claim.description,
      type: 'swarm_claim',
      timestamp: Date.now(),
      details: []
    };

    // Based on claim type, perform appropriate verifications
    switch (claim.type) {
      case 'function_implementation':
        // Verify required types exist
        for (const type of claim.requiredTypes || []) {
          const result = this.verifyTypeExists(claim.packagePath, type);
          verification.details.push(result);
        }
        
        // Verify required functions exist
        for (const func of claim.requiredFunctions || []) {
          const result = this.verifyFunctionExists(claim.packagePath, func.name, func.receiver);
          verification.details.push(result);
        }
        break;

      case 'type_creation':
        // Verify the type doesn't already exist (to avoid duplicates)
        const existing = this.verifyTypeExists(claim.packagePath, claim.typeName);
        verification.details.push({
          ...existing,
          claim: `Type ${claim.typeName} should not already exist`,
          verified: !existing.verified // Invert - we want it to NOT exist
        });
        break;

      case 'import_fix':
        // Verify the import is actually needed
        const usage = this.verifyImportUsed(claim.filePath, claim.importPath);
        verification.details.push(usage);
        break;
    }

    // Calculate overall verification
    verification.verified = verification.details.length > 0 && 
      verification.details.every(d => d.verified);

    this.recordVerification(verification);
    return verification;
  }

  /**
   * Verify compilation will work
   */
  async verifyCompilation(packagePath) {
    const verification = {
      claim: `Package ${packagePath} compiles`,
      type: 'compilation',
      verified: false,
      evidence: {}
    };

    try {
      // Try to build the package
      const goPath = process.env.GOPATH || path.join(process.env.HOME, 'go');
      const result = execSync(`cd ${packagePath} && ${goPath}/bin/go build -o /dev/null ./...`, {
        encoding: 'utf8',
        stdio: 'pipe'
      });
      
      verification.verified = true;
      verification.evidence = {
        compiled: true,
        output: result
      };
    } catch (error) {
      verification.evidence = {
        compiled: false,
        errors: error.stderr || error.message
      };
      
      // Parse specific errors
      if (error.stderr) {
        const errors = this.parseGoErrors(error.stderr);
        verification.evidence.parsedErrors = errors;
      }
    }

    this.recordVerification(verification);
    return verification;
  }

  /**
   * Parse Go compilation errors
   */
  parseGoErrors(stderr) {
    const errors = [];
    const lines = stderr.split('\n');
    
    for (const line of lines) {
      // Parse errors like: x/module/file.go:10:5: undefined: SomeType
      const match = line.match(/([^:]+):(\d+):(\d+):\s*(.+)/);
      if (match) {
        errors.push({
          file: match[1],
          line: parseInt(match[2]),
          column: parseInt(match[3]),
          message: match[4]
        });
      }
    }
    
    return errors;
  }

  /**
   * Find all Go files in a directory
   */
  findGoFiles(dirPath) {
    const files = [];
    
    try {
      const walkDir = (dir) => {
        const entries = fs.readdirSync(dir, { withFileTypes: true });
        for (const entry of entries) {
          const fullPath = path.join(dir, entry.name);
          if (entry.isDirectory() && !entry.name.startsWith('.')) {
            walkDir(fullPath);
          } else if (entry.name.endsWith('.go') && !entry.name.endsWith('_test.go')) {
            files.push(fullPath);
          }
        }
      };
      walkDir(dirPath);
    } catch (error) {
      // Directory might not exist
    }
    
    return files;
  }

  /**
   * Record a verification
   */
  recordVerification(verification) {
    verification.agent = this.name;
    verification.timestamp = verification.timestamp || Date.now();
    
    this.verifications.push(verification);
    this.witness('VERIFICATION', verification);
    
    // Update veracity score
    if (verification.verified) {
      console.log(`✅ VERIFIED: ${verification.claim}`);
    } else {
      console.log(`❌ FALSE: ${verification.claim}`);
      if (verification.evidence) {
        console.log(`   Evidence:`, verification.evidence);
      }
    }
  }

  /**
   * Generate veracity report
   */
  generateReport() {
    const report = {
      agent: this.name,
      lctId: this.lctId,
      timestamp: Date.now(),
      totalVerifications: this.verifications.length,
      verified: this.verifications.filter(v => v.verified).length,
      failed: this.verifications.filter(v => !v.verified).length,
      veracityScore: 0,
      breakdownByType: {}
    };

    // Calculate veracity score
    if (this.verifications.length > 0) {
      report.veracityScore = report.verified / this.verifications.length;
    }

    // Breakdown by type
    for (const v of this.verifications) {
      const type = v.type || 'unknown';
      if (!report.breakdownByType[type]) {
        report.breakdownByType[type] = { total: 0, verified: 0 };
      }
      report.breakdownByType[type].total++;
      if (v.verified) {
        report.breakdownByType[type].verified++;
      }
    }

    return report;
  }

  /**
   * Witness an event
   */
  witness(action, details) {
    const entry = {
      timestamp: Date.now(),
      agent: this.name,
      action,
      ...details
    };
    
    // Ensure directory exists
    const dir = path.dirname(this.witnessLog);
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }

    fs.appendFileSync(this.witnessLog, JSON.stringify(entry) + '\n');
  }
}

// Export for use in swarm
module.exports = VeracityAgent;

// CLI interface for testing
if (require.main === module) {
  const veracity = new VeracityAgent();
  
  // Example verifications
  const ledgerPath = path.join(__dirname, '../implementation/ledger');
  
  console.log('🔍 Running Veracity Checks...\n');
  
  // Check if claimed types exist
  const checks = [
    { package: path.join(ledgerPath, 'x/lctmanager/types'), type: 'LCT' },
    { package: path.join(ledgerPath, 'x/lctmanager/types'), type: 'MRH' },
    { package: path.join(ledgerPath, 'x/lctmanager/types'), type: 'BirthCertificate' },
    { package: path.join(ledgerPath, 'x/energycycle/types'), type: 'EnergyPool' },
    { package: path.join(ledgerPath, 'x/energycycle/types'), type: 'ATPToken' },
    { package: path.join(ledgerPath, 'x/energycycle/types'), type: 'ADPToken' },
    { package: path.join(ledgerPath, 'x/trusttensor/types'), type: 'TrustRecord' },
  ];
  
  for (const check of checks) {
    veracity.verifyTypeExists(check.package, check.type);
  }
  
  // Check if claimed functions exist
  const funcChecks = [
    { package: path.join(ledgerPath, 'x/lctmanager/keeper'), func: 'MintLCT', receiver: 'Keeper' },
    { package: path.join(ledgerPath, 'x/lctmanager/keeper'), func: 'BindLCT', receiver: 'Keeper' },
    { package: path.join(ledgerPath, 'x/energycycle/keeper'), func: 'DischargeATP', receiver: 'Keeper' },
    { package: path.join(ledgerPath, 'x/trusttensor/keeper'), func: 'UpdateT3', receiver: 'Keeper' },
  ];
  
  for (const check of funcChecks) {
    veracity.verifyFunctionExists(check.package, check.func, check.receiver);
  }
  
  // Check if fields exist
  const fieldChecks = [
    { package: path.join(ledgerPath, 'x/lctmanager/keeper'), struct: 'Keeper', field: 'storeKey' },
    { package: path.join(ledgerPath, 'x/lctmanager/types'), struct: 'LCT', field: 'T3Tensor' },
  ];
  
  for (const check of fieldChecks) {
    veracity.verifyFieldExists(check.package, check.struct, check.field);
  }
  
  console.log('\n📊 Veracity Report:');
  console.log(JSON.stringify(veracity.generateReport(), null, 2));
}