#!/usr/bin/env node

/**
 * Swarm Validator Agent
 * Responsible for verifying that claimed work actually happened
 * Updates V3 (Veracity, Validity, Value) tensors based on validation results
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

class ValidatorAgent {
  constructor(config = {}) {
    this.name = config.name || 'Validator-Agent';
    this.lctId = `lct:validator:${Date.now()}`;
    this.witnessLog = config.witnessLog || './swarm-memory/witness/validator.jsonl';
    this.v3Updates = [];
  }

  /**
   * Validate claimed file creation
   */
  validateFileExists(filePath, expectedContent = null) {
    const result = {
      claim: `File exists: ${filePath}`,
      valid: false,
      evidence: null
    };

    try {
      if (fs.existsSync(filePath)) {
        result.valid = true;
        result.evidence = {
          exists: true,
          size: fs.statSync(filePath).size,
          modified: fs.statSync(filePath).mtime
        };

        if (expectedContent) {
          const actualContent = fs.readFileSync(filePath, 'utf8');
          result.valid = actualContent.includes(expectedContent);
          result.evidence.contentMatch = result.valid;
        }
      }
    } catch (error) {
      result.evidence = { error: error.message };
    }

    this.updateV3(result);
    return result;
  }

  /**
   * Validate claimed command execution
   */
  validateCommandRuns(command, expectedOutput = null) {
    const result = {
      claim: `Command runs: ${command}`,
      valid: false,
      evidence: null
    };

    try {
      const output = execSync(command, { encoding: 'utf8', timeout: 5000 });
      result.valid = true;
      result.evidence = {
        executed: true,
        outputLength: output.length,
        exitCode: 0
      };

      if (expectedOutput) {
        result.valid = output.includes(expectedOutput);
        result.evidence.outputMatch = result.valid;
      }
    } catch (error) {
      result.evidence = {
        executed: false,
        error: error.message,
        exitCode: error.status
      };
    }

    this.updateV3(result);
    return result;
  }

  /**
   * Validate claimed compilation
   */
  validateCompiles(sourcePath, binaryPath = null) {
    const result = {
      claim: `Code compiles: ${sourcePath}`,
      valid: false,
      evidence: null
    };

    try {
      // Check if source exists
      if (!fs.existsSync(sourcePath)) {
        result.evidence = { sourceExists: false };
        this.updateV3(result);
        return result;
      }

      // Try to compile based on file extension
      const ext = path.extname(sourcePath);
      let compileCmd = null;

      switch(ext) {
        case '.go':
          compileCmd = `go build -o /tmp/test_binary ${sourcePath}`;
          break;
        case '.js':
          compileCmd = `node -c ${sourcePath}`;
          break;
        case '.py':
          compileCmd = `python3 -m py_compile ${sourcePath}`;
          break;
        default:
          result.evidence = { error: 'Unknown file type' };
          this.updateV3(result);
          return result;
      }

      execSync(compileCmd, { timeout: 10000 });
      result.valid = true;
      result.evidence = { compiled: true };

      // Check if binary was created
      if (binaryPath && fs.existsSync(binaryPath)) {
        result.evidence.binaryCreated = true;
        result.evidence.binarySize = fs.statSync(binaryPath).size;
      }

    } catch (error) {
      result.evidence = {
        compiled: false,
        error: error.message
      };
    }

    this.updateV3(result);
    return result;
  }

  /**
   * Validate blockchain specific claims
   */
  validateBlockchainClaim(claim) {
    const validations = {
      'proto_generated': () => this.validateFileExists('./x/lctmanager/types/lct.pb.go'),
      'keeper_implemented': () => this.validateFileExists('./x/lctmanager/keeper/lct_lifecycle.go', 'MintLCT'),
      'module_wired': () => this.validateFileExists('./app/app.go', 'LCTManagerKeeper'),
      'genesis_valid': () => this.validateJsonFile('./config/genesis.json'),
      'chain_runs': () => this.validateCommandRuns('./actd version'),
      'lct_mintable': () => this.validateCommandRuns('./actd tx lctmanager mint-lct', '--help'),
      'atp_functional': () => this.validateCommandRuns('./actd query energycycle pool', 'default')
    };

    if (validations[claim]) {
      return validations[claim]();
    }

    return {
      claim: claim,
      valid: false,
      evidence: { error: 'Unknown claim type' }
    };
  }

  /**
   * Validate JSON file structure
   */
  validateJsonFile(filePath) {
    const result = {
      claim: `Valid JSON: ${filePath}`,
      valid: false,
      evidence: null
    };

    try {
      const content = fs.readFileSync(filePath, 'utf8');
      const parsed = JSON.parse(content);
      result.valid = true;
      result.evidence = {
        valid: true,
        keys: Object.keys(parsed).length
      };
    } catch (error) {
      result.evidence = {
        valid: false,
        error: error.message
      };
    }

    this.updateV3(result);
    return result;
  }

  /**
   * Update V3 tensor based on validation result
   */
  updateV3(validationResult) {
    const v3Update = {
      timestamp: Date.now(),
      claim: validationResult.claim,
      veracity: validationResult.valid ? 1.0 : 0.0,  // Truth of claim
      validity: validationResult.evidence ? 0.8 : 0.2, // Quality of evidence
      value: validationResult.valid ? 0.7 : 0.3      // Value delivered
    };

    this.v3Updates.push(v3Update);
    this.witness('V3_UPDATE', v3Update);
    return v3Update;
  }

  /**
   * Validate an entire swarm execution
   */
  async validateSwarmExecution(swarmReport) {
    const validation = {
      swarmName: swarmReport.name,
      timestamp: Date.now(),
      claims: [],
      overallV3: { veracity: 0, validity: 0, value: 0 }
    };

    // Validate each claimed deliverable
    for (const deliverable of swarmReport.deliverables || []) {
      const result = await this.validateDeliverable(deliverable);
      validation.claims.push(result);
    }

    // Calculate overall V3 scores
    if (validation.claims.length > 0) {
      validation.overallV3.veracity = validation.claims.reduce((sum, c) => sum + (c.valid ? 1 : 0), 0) / validation.claims.length;
      validation.overallV3.validity = validation.claims.reduce((sum, c) => sum + (c.evidence ? 0.8 : 0.2), 0) / validation.claims.length;
      validation.overallV3.value = validation.claims.reduce((sum, c) => sum + (c.valid ? c.value || 0.5 : 0), 0) / validation.claims.length;
    }

    this.witness('SWARM_VALIDATION_COMPLETE', validation);
    return validation;
  }

  /**
   * Validate individual deliverable
   */
  async validateDeliverable(deliverable) {
    const validators = {
      'file': () => this.validateFileExists(deliverable.path, deliverable.contains),
      'command': () => this.validateCommandRuns(deliverable.command, deliverable.expectedOutput),
      'compilation': () => this.validateCompiles(deliverable.source, deliverable.binary),
      'blockchain': () => this.validateBlockchainClaim(deliverable.claim),
      'json': () => this.validateJsonFile(deliverable.path)
    };

    const validator = validators[deliverable.type];
    if (validator) {
      const result = validator();
      result.deliverable = deliverable.name;
      result.value = deliverable.value || 0.5;
      return result;
    }

    return {
      deliverable: deliverable.name,
      valid: false,
      evidence: { error: 'Unknown deliverable type' }
    };
  }

  /**
   * Generate validation report
   */
  generateReport() {
    const report = {
      validator: this.name,
      lctId: this.lctId,
      timestamp: Date.now(),
      validations: this.v3Updates.length,
      averageV3: {
        veracity: 0,
        validity: 0,
        value: 0
      }
    };

    if (this.v3Updates.length > 0) {
      report.averageV3.veracity = this.v3Updates.reduce((sum, u) => sum + u.veracity, 0) / this.v3Updates.length;
      report.averageV3.validity = this.v3Updates.reduce((sum, u) => sum + u.validity, 0) / this.v3Updates.length;
      report.averageV3.value = this.v3Updates.reduce((sum, u) => sum + u.value, 0) / this.v3Updates.length;
    }

    return report;
  }

  /**
   * Witness an event
   */
  witness(action, details) {
    const entry = {
      timestamp: Date.now(),
      validator: this.name,
      action,
      ...details
    };
    
    // Ensure directory exists
    const dir = path.dirname(this.witnessLog);
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }

    fs.appendFileSync(this.witnessLog, JSON.stringify(entry) + '\n');
    console.log(`👁️ [VALIDATOR] ${action}:`, details.claim || details.swarmName || '');
  }
}

// Export for use in swarm
module.exports = ValidatorAgent;

// CLI interface
if (require.main === module) {
  const validator = new ValidatorAgent();
  
  // Example validation of blockchain swarm deliverables
  const blockchainDeliverables = [
    { type: 'file', name: 'LCT Keeper', path: './x/lctmanager/keeper/lct_lifecycle.go', contains: 'MintLCT', value: 0.8 },
    { type: 'file', name: 'ATP/ADP System', path: './x/energycycle/keeper/atp_adp.go', contains: 'DischargeATP', value: 0.9 },
    { type: 'file', name: 'Trust Attribution', path: './x/trusttensor/keeper/trust_attribution.go', contains: 'UpdateT3', value: 0.7 },
    { type: 'json', name: 'Genesis Config', path: './config/genesis.json', value: 0.6 },
    { type: 'blockchain', name: 'Proto Generation', claim: 'proto_generated', value: 1.0 },
    { type: 'blockchain', name: 'Chain Runs', claim: 'chain_runs', value: 1.0 }
  ];

  console.log('🔍 Validating Blockchain Swarm Deliverables...\n');
  
  blockchainDeliverables.forEach(deliverable => {
    const result = validator.validateDeliverable(deliverable);
    const symbol = result.valid ? '✅' : '❌';
    console.log(`${symbol} ${deliverable.name}: ${result.valid ? 'VALID' : 'INVALID'}`);
    if (result.evidence) {
      console.log(`   Evidence:`, result.evidence);
    }
  });

  console.log('\n📊 Validation Report:');
  console.log(JSON.stringify(validator.generateReport(), null, 2));
}