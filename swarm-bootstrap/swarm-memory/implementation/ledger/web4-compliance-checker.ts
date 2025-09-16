/**
 * Web4 Compliance Checker Role
 * A specialized role called by the Law Oracle to verify Web4 standard compliance
 * 
 * This role validates:
 * - Proper LCT structure and binding
 * - R6 rule compliance
 * - MRH boundaries
 * - Witnessed presence protocol
 * - ATP/ADP economics
 * - Society governance rules
 */

import { Web4Transaction, Web4Block, TransactionType } from './web4-blockchain';

// Compliance check result
export interface ComplianceResult {
  compliant: boolean;
  score: number;        // 0-100 compliance score
  violations: ComplianceViolation[];
  suggestions: string[];
  timestamp: number;
  checkerLCT: string;
  signature: string;
}

export interface ComplianceViolation {
  severity: 'critical' | 'high' | 'medium' | 'low';
  category: ComplianceCategory;
  description: string;
  specification: string;  // Which part of Web4 spec violated
  remedy?: string;        // How to fix
}

export enum ComplianceCategory {
  // Core Web4 violations
  LCT_STRUCTURE = 'lct_structure',
  BINDING_PROTOCOL = 'binding_protocol',
  WITNESS_PRESENCE = 'witness_presence',
  MRH_BOUNDARIES = 'mrh_boundaries',
  
  // R6 violations
  R6_RULES = 'r6_rules',
  R6_ROLES = 'r6_roles',
  R6_REQUEST = 'r6_request',
  R6_REFERENCE = 'r6_reference',
  R6_RESOURCE = 'r6_resource',
  R6_RESULT = 'r6_result',
  
  // Economic violations
  ATP_BALANCE = 'atp_balance',
  ADP_GENERATION = 'adp_generation',
  ECONOMIC_MODEL = 'economic_model',
  
  // Governance violations
  SOCIETY_RULES = 'society_rules',
  LAW_ORACLE = 'law_oracle',
  CITIZEN_RIGHTS = 'citizen_rights',
  
  // Technical violations
  CRYPTOGRAPHY = 'cryptography',
  SIGNATURE = 'signature',
  TIMESTAMP = 'timestamp'
}

// R6 Compliance Rules
interface R6Rules {
  Rules: string[];
  Roles: string[];
  Request: string[];
  Reference: string[];
  Resource: string[];
  Result: string[];
}

/**
 * Web4 Compliance Checker Role Implementation
 */
export class Web4ComplianceChecker {
  private roleLCT: string;
  private lawOracleLCT: string;
  private society: string;
  private r6Rules: R6Rules;
  
  constructor(roleLCT: string, lawOracleLCT: string, society: string) {
    this.roleLCT = roleLCT;
    this.lawOracleLCT = lawOracleLCT;
    this.society = society;
    
    // Define R6 rules for compliance checker role
    this.r6Rules = {
      Rules: [
        'Validate all Web4 specifications',
        'Ensure cryptographic integrity',
        'Verify witness presence',
        'Check R6 compliance'
      ],
      Roles: [
        'Compliance validator',
        'Standard enforcer',
        'Law oracle assistant'
      ],
      Request: [
        'Check transaction compliance',
        'Validate block integrity',
        'Audit entity behavior'
      ],
      Reference: [
        'Web4 Protocol Specification',
        'LCT Binding Standard',
        'R6 Framework',
        'Society Law Oracle'
      ],
      Resource: [
        '5 ATP per check',
        'Access to blockchain',
        'Law oracle authority'
      ],
      Result: [
        'Compliance report',
        'Violation list',
        'Remediation suggestions'
      ]
    };
  }
  
  /**
   * Check transaction compliance
   */
  async checkTransaction(transaction: Web4Transaction): Promise<ComplianceResult> {
    const violations: ComplianceViolation[] = [];
    
    // 1. Check LCT structure
    if (!this.isValidLCT(transaction.from)) {
      violations.push({
        severity: 'critical',
        category: ComplianceCategory.LCT_STRUCTURE,
        description: `Invalid LCT format: ${transaction.from}`,
        specification: 'Web4 Spec 2.1.3 - LCT Structure',
        remedy: 'LCT must follow format: lct:web4:entity_type:identifier'
      });
    }
    
    // 2. Check signature
    if (!this.isValidSignature(transaction)) {
      violations.push({
        severity: 'critical',
        category: ComplianceCategory.SIGNATURE,
        description: 'Invalid or missing transaction signature',
        specification: 'Web4 Spec 3.2 - Cryptographic Binding',
        remedy: 'Sign transaction with Ed25519 private key'
      });
    }
    
    // 3. Check witness requirements
    const witnessRequired = this.witnessesRequired(transaction.type);
    if (transaction.witnesses.length < witnessRequired) {
      violations.push({
        severity: 'high',
        category: ComplianceCategory.WITNESS_PRESENCE,
        description: `Insufficient witnesses: ${transaction.witnesses.length}/${witnessRequired}`,
        specification: 'Web4 Spec 4.1 - Witnessed Presence Protocol',
        remedy: `Obtain ${witnessRequired} witness signatures`
      });
    }
    
    // 4. Check type-specific compliance
    const typeViolations = await this.checkTransactionType(transaction);
    violations.push(...typeViolations);
    
    // 5. Check R6 compliance if role transaction
    if (transaction.type.startsWith('role_')) {
      const r6Violations = this.checkR6Compliance(transaction);
      violations.push(...r6Violations);
    }
    
    // 6. Check ATP/ADP economics
    if (transaction.type.startsWith('atp_') || transaction.type.startsWith('adp_')) {
      const economicViolations = this.checkEconomicCompliance(transaction);
      violations.push(...economicViolations);
    }
    
    // Calculate compliance score
    const score = this.calculateComplianceScore(violations);
    
    // Generate suggestions
    const suggestions = this.generateSuggestions(violations);
    
    return {
      compliant: violations.filter(v => v.severity === 'critical').length === 0,
      score,
      violations,
      suggestions,
      timestamp: Date.now(),
      checkerLCT: this.roleLCT,
      signature: await this.signResult(violations)
    };
  }
  
  /**
   * Check block compliance
   */
  async checkBlock(block: Web4Block): Promise<ComplianceResult> {
    const violations: ComplianceViolation[] = [];
    
    // 1. Check block structure
    if (!block.society || !block.lawOracle) {
      violations.push({
        severity: 'critical',
        category: ComplianceCategory.SOCIETY_RULES,
        description: 'Block missing society or law oracle',
        specification: 'Web4 Spec 5.1 - Society Governance',
        remedy: 'Blocks must specify owning society and law oracle'
      });
    }
    
    // 2. Check witness presence
    if (block.witnessedBy.length === 0) {
      violations.push({
        severity: 'critical',
        category: ComplianceCategory.WITNESS_PRESENCE,
        description: 'Block has no witnesses',
        specification: 'Web4 Spec 4.2 - Block Witnessing',
        remedy: 'At least 3 witnesses required for block validity'
      });
    }
    
    // 3. Check MRH boundaries
    const mrhViolations = this.checkMRHBoundaries(block.mrh);
    violations.push(...mrhViolations);
    
    // 4. Check all transactions in block
    for (const tx of block.transactions) {
      const txResult = await this.checkTransaction(tx);
      if (!txResult.compliant) {
        violations.push({
          severity: 'high',
          category: ComplianceCategory.LCT_STRUCTURE,
          description: `Transaction ${tx.id} is non-compliant`,
          specification: 'Web4 Spec 6.1 - Transaction Validity',
          remedy: 'All transactions in block must be compliant'
        });
      }
    }
    
    // 5. Check cryptographic integrity
    if (!this.verifyBlockHash(block)) {
      violations.push({
        severity: 'critical',
        category: ComplianceCategory.CRYPTOGRAPHY,
        description: 'Block hash verification failed',
        specification: 'Web4 Spec 3.4 - Block Integrity',
        remedy: 'Recalculate block hash with proper nonce'
      });
    }
    
    const score = this.calculateComplianceScore(violations);
    const suggestions = this.generateSuggestions(violations);
    
    return {
      compliant: violations.filter(v => v.severity === 'critical').length === 0,
      score,
      violations,
      suggestions,
      timestamp: Date.now(),
      checkerLCT: this.roleLCT,
      signature: await this.signResult(violations)
    };
  }
  
  /**
   * Validate LCT format
   */
  private isValidLCT(lct: string): boolean {
    const lctPattern = /^lct:web4:(human|ai|role|society|dictionary):[a-zA-Z0-9-]+$/;
    return lctPattern.test(lct);
  }
  
  /**
   * Validate signature (simplified)
   */
  private isValidSignature(transaction: Web4Transaction): boolean {
    // In production, would verify Ed25519 signature
    return transaction.signature && transaction.signature.length > 0;
  }
  
  /**
   * Determine witness requirements by transaction type
   */
  private witnessesRequired(type: TransactionType): number {
    switch (type) {
      case TransactionType.SOCIETY_CREATE:
        return 5;  // High witness requirement for society creation
      case TransactionType.LCT_CREATE:
        return 3;  // Moderate for LCT creation
      case TransactionType.ATP_TRANSFER:
        return 2;  // Lower for simple transfers
      case TransactionType.WITNESS_ACTION:
        return 1;  // Minimal for witness actions
      default:
        return 2;  // Default requirement
    }
  }
  
  /**
   * Check transaction type specific rules
   */
  private async checkTransactionType(tx: Web4Transaction): Promise<ComplianceViolation[]> {
    const violations: ComplianceViolation[] = [];
    
    switch (tx.type) {
      case TransactionType.LCT_CREATE:
        if (!tx.data.lct?.publicKey) {
          violations.push({
            severity: 'critical',
            category: ComplianceCategory.BINDING_PROTOCOL,
            description: 'LCT creation missing public key',
            specification: 'Web4 Spec 2.2 - Cryptographic Binding',
            remedy: 'Provide Ed25519 public key for binding'
          });
        }
        break;
        
      case TransactionType.ATP_TRANSFER:
        if (!tx.data.atp?.amount || tx.data.atp.amount <= 0) {
          violations.push({
            severity: 'high',
            category: ComplianceCategory.ATP_BALANCE,
            description: 'Invalid ATP transfer amount',
            specification: 'Web4 Spec 7.1 - ATP Transfers',
            remedy: 'Amount must be positive integer'
          });
        }
        break;
        
      case TransactionType.ROLE_EXECUTE:
        if (!tx.data.role?.r6Rules) {
          violations.push({
            severity: 'high',
            category: ComplianceCategory.R6_RULES,
            description: 'Role execution missing R6 rules',
            specification: 'Web4 Spec 8.1 - Role Framework',
            remedy: 'Provide complete R6 rule set'
          });
        }
        break;
    }
    
    return violations;
  }
  
  /**
   * Check R6 compliance for role transactions
   */
  private checkR6Compliance(tx: Web4Transaction): ComplianceViolation[] {
    const violations: ComplianceViolation[] = [];
    
    if (!tx.data.role?.r6Rules) {
      return violations;
    }
    
    const r6 = tx.data.role.r6Rules;
    
    // Check all R6 components present
    const components = ['Rules', 'Roles', 'Request', 'Reference', 'Resource', 'Result'];
    for (const component of components) {
      if (!r6[component] || r6[component].length === 0) {
        violations.push({
          severity: 'high',
          category: ComplianceCategory[`R6_${component.toUpperCase()}`],
          description: `R6 ${component} not defined`,
          specification: `Web4 Spec 8.2 - R6 ${component}`,
          remedy: `Define ${component} in R6 framework`
        });
      }
    }
    
    return violations;
  }
  
  /**
   * Check economic compliance
   */
  private checkEconomicCompliance(tx: Web4Transaction): ComplianceViolation[] {
    const violations: ComplianceViolation[] = [];
    
    if (tx.type === TransactionType.ADP_GENERATE) {
      if (!tx.data.adp?.r6Compliance) {
        violations.push({
          severity: 'high',
          category: ComplianceCategory.ADP_GENERATION,
          description: 'ADP generation missing R6 compliance proof',
          specification: 'Web4 Spec 7.3 - ADP Generation',
          remedy: 'Provide R6 compliance proof for ADP'
        });
      }
      
      // Check ADP generation ratio
      if (tx.data.adp?.atpConsumed && tx.data.adp?.adpGenerated) {
        const ratio = tx.data.adp.adpGenerated / tx.data.adp.atpConsumed;
        if (ratio > 2) {
          violations.push({
            severity: 'medium',
            category: ComplianceCategory.ECONOMIC_MODEL,
            description: `ADP generation ratio too high: ${ratio}`,
            specification: 'Web4 Spec 7.4 - Economic Balance',
            remedy: 'ADP generation should not exceed 2x ATP consumed'
          });
        }
      }
    }
    
    return violations;
  }
  
  /**
   * Check MRH boundaries
   */
  private checkMRHBoundaries(mrh: any): ComplianceViolation[] {
    const violations: ComplianceViolation[] = [];
    
    if (!mrh || !mrh.bound || !mrh.witnessing) {
      violations.push({
        severity: 'high',
        category: ComplianceCategory.MRH_BOUNDARIES,
        description: 'Missing MRH boundary definitions',
        specification: 'Web4 Spec 9.1 - Markov Relevancy Horizon',
        remedy: 'Define bound and witnessing entities'
      });
    }
    
    // Check for orphaned entities (in paired but not bound)
    if (mrh?.paired) {
      for (const entity of mrh.paired) {
        if (!mrh.bound.includes(entity)) {
          violations.push({
            severity: 'medium',
            category: ComplianceCategory.MRH_BOUNDARIES,
            description: `Entity ${entity} paired but not bound`,
            specification: 'Web4 Spec 9.2 - MRH Consistency',
            remedy: 'Paired entities must also be bound'
          });
        }
      }
    }
    
    return violations;
  }
  
  /**
   * Verify block hash
   */
  private verifyBlockHash(block: Web4Block): boolean {
    // In production, would recalculate and verify hash
    return block.hash && block.hash.length === 64;
  }
  
  /**
   * Calculate compliance score
   */
  private calculateComplianceScore(violations: ComplianceViolation[]): number {
    if (violations.length === 0) return 100;
    
    let score = 100;
    
    for (const violation of violations) {
      switch (violation.severity) {
        case 'critical':
          score -= 25;
          break;
        case 'high':
          score -= 15;
          break;
        case 'medium':
          score -= 10;
          break;
        case 'low':
          score -= 5;
          break;
      }
    }
    
    return Math.max(0, score);
  }
  
  /**
   * Generate remediation suggestions
   */
  private generateSuggestions(violations: ComplianceViolation[]): string[] {
    const suggestions: string[] = [];
    
    // Group violations by category
    const byCategory = new Map<ComplianceCategory, ComplianceViolation[]>();
    for (const v of violations) {
      if (!byCategory.has(v.category)) {
        byCategory.set(v.category, []);
      }
      byCategory.get(v.category)!.push(v);
    }
    
    // Generate category-specific suggestions
    if (byCategory.has(ComplianceCategory.WITNESS_PRESENCE)) {
      suggestions.push('Increase witness participation by incentivizing with ATP rewards');
    }
    
    if (byCategory.has(ComplianceCategory.R6_RULES)) {
      suggestions.push('Review and complete R6 framework definitions for all roles');
    }
    
    if (byCategory.has(ComplianceCategory.ATP_BALANCE)) {
      suggestions.push('Implement balance checking before ATP transactions');
    }
    
    if (violations.filter(v => v.severity === 'critical').length > 0) {
      suggestions.unshift('URGENT: Address critical violations before proceeding');
    }
    
    return suggestions;
  }
  
  /**
   * Sign compliance result
   */
  private async signResult(violations: ComplianceViolation[]): Promise<string> {
    // In production, would sign with Ed25519 private key
    const data = JSON.stringify(violations);
    return `signature:${Buffer.from(data).toString('base64').substring(0, 16)}`;
  }
  
  /**
   * Create compliance report
   */
  async createComplianceReport(
    entity: string,
    transactions: Web4Transaction[],
    blocks: Web4Block[]
  ): Promise<{
    entity: string;
    overallCompliance: number;
    transactionCompliance: ComplianceResult[];
    blockCompliance: ComplianceResult[];
    summary: string;
  }> {
    const txResults: ComplianceResult[] = [];
    const blockResults: ComplianceResult[] = [];
    
    // Check all transactions
    for (const tx of transactions) {
      const result = await this.checkTransaction(tx);
      txResults.push(result);
    }
    
    // Check all blocks
    for (const block of blocks) {
      const result = await this.checkBlock(block);
      blockResults.push(result);
    }
    
    // Calculate overall compliance
    const allScores = [...txResults, ...blockResults].map(r => r.score);
    const overallCompliance = allScores.reduce((a, b) => a + b, 0) / allScores.length;
    
    // Generate summary
    const criticalViolations = [...txResults, ...blockResults]
      .flatMap(r => r.violations)
      .filter(v => v.severity === 'critical');
    
    const summary = criticalViolations.length > 0
      ? `Entity has ${criticalViolations.length} critical violations requiring immediate attention`
      : `Entity is Web4 compliant with ${overallCompliance.toFixed(1)}% compliance score`;
    
    return {
      entity,
      overallCompliance,
      transactionCompliance: txResults,
      blockCompliance: blockResults,
      summary
    };
  }
}

// Export for use in law oracle
export default Web4ComplianceChecker;