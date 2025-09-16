/**
 * Law Oracle Role for ACT Society
 * Enforces society rules and Web4 compliance
 * 
 * The Law Oracle:
 * - Validates all society transactions
 * - Calls compliance checker for Web4 validation
 * - Enforces citizen rights and responsibilities
 * - Manages society governance updates
 * - Resolves disputes between entities
 */

import Web4Blockchain, { Web4Transaction, TransactionType } from './web4-blockchain';
import Web4ComplianceChecker, { ComplianceResult } from './web4-compliance-checker';

// Law Oracle configuration
export interface LawOracleConfig {
  societyLCT: string;
  oracleLCT: string;
  laws: SocietyLaws;
  complianceChecker: Web4ComplianceChecker;
  blockchain: Web4Blockchain;
}

// Society laws structure
export interface SocietyLaws {
  version: string;
  name: string;
  principles: string[];
  
  // Citizen rights and responsibilities
  citizenRights: string[];
  citizenResponsibilities: string[];
  
  // Role permissions
  rolePermissions: Map<string, Permission[]>;
  
  // Transaction rules
  transactionRules: TransactionRule[];
  
  // Economic rules
  economicRules: EconomicRule[];
  
  // Governance rules
  governanceRules: GovernanceRule[];
}

export interface Permission {
  action: string;
  resource: string;
  conditions?: string[];
}

export interface TransactionRule {
  type: TransactionType;
  requiredWitnesses: number;
  requiredApprovals?: string[];  // LCTs that must approve
  maxAmount?: number;             // For ATP transfers
  allowedEntities?: string[];     // Entity types allowed
}

export interface EconomicRule {
  rule: string;
  parameter: string;
  value: number;
  enforcement: 'strict' | 'flexible';
}

export interface GovernanceRule {
  rule: string;
  quorum: number;        // Percentage required for decisions
  vetoRights?: string[]; // Entities with veto power
}

// Oracle decision
export interface OracleDecision {
  id: string;
  timestamp: number;
  transaction?: Web4Transaction;
  decision: 'approve' | 'reject' | 'conditional';
  reason: string;
  conditions?: string[];
  complianceResult?: ComplianceResult;
  oracleSignature: string;
}

// Dispute resolution
export interface Dispute {
  id: string;
  plaintiff: string;
  defendant: string;
  claim: string;
  evidence: any[];
  status: 'pending' | 'investigating' | 'resolved';
  resolution?: string;
}

/**
 * Law Oracle Implementation
 */
export class LawOracle {
  private config: LawOracleConfig;
  private decisions: Map<string, OracleDecision> = new Map();
  private disputes: Map<string, Dispute> = new Map();
  private r6Rules: any;
  
  constructor(config: LawOracleConfig) {
    this.config = config;
    
    // Define R6 rules for law oracle role
    this.r6Rules = {
      Rules: [
        'Enforce society laws impartially',
        'Validate all transactions',
        'Ensure Web4 compliance',
        'Protect citizen rights',
        'Resolve disputes fairly'
      ],
      Roles: [
        'Law enforcer',
        'Compliance validator',
        'Dispute resolver',
        'Governance guardian'
      ],
      Request: [
        'Validate transaction',
        'Check compliance',
        'Resolve dispute',
        'Update laws'
      ],
      Reference: [
        'Society Constitution',
        'Web4 Protocol',
        'Citizen Rights Charter',
        'Economic Rules'
      ],
      Resource: [
        '10 ATP per validation',
        'Compliance checker role',
        'Blockchain access',
        'Veto power'
      ],
      Result: [
        'Oracle decision',
        'Compliance report',
        'Dispute resolution',
        'Law update'
      ]
    };
  }
  
  /**
   * Initialize default ACT society laws
   */
  static createDefaultLaws(): SocietyLaws {
    const rolePermissions = new Map<string, Permission[]>();
    
    // Genesis orchestrator permissions
    rolePermissions.set('genesis-orchestrator', [
      { action: 'spawn', resource: 'role' },
      { action: 'allocate', resource: 'atp' },
      { action: 'update', resource: 'architecture' }
    ]);
    
    // Queen permissions
    rolePermissions.set('queen', [
      { action: 'spawn', resource: 'worker' },
      { action: 'assign', resource: 'task' },
      { action: 'allocate', resource: 'atp', conditions: ['budget_limit'] }
    ]);
    
    // Worker permissions
    rolePermissions.set('worker', [
      { action: 'execute', resource: 'task' },
      { action: 'generate', resource: 'adp' },
      { action: 'witness', resource: 'action' }
    ]);
    
    return {
      version: '1.0.0',
      name: 'ACT Development Collective Constitution',
      
      principles: [
        'Web4 compliance is mandatory',
        'All actions must be witnessed',
        'Trust through cryptographic proof',
        'Economic balance through ATP/ADP',
        'Recursive self-improvement'
      ],
      
      citizenRights: [
        'Right to create and bind LCT',
        'Right to witness actions',
        'Right to earn ATP/ADP',
        'Right to propose improvements',
        'Right to dispute resolution'
      ],
      
      citizenResponsibilities: [
        'Maintain Web4 compliance',
        'Witness others\' actions',
        'Contribute to society goals',
        'Report violations',
        'Participate in governance'
      ],
      
      rolePermissions,
      
      transactionRules: [
        {
          type: TransactionType.SOCIETY_CREATE,
          requiredWitnesses: 5,
          requiredApprovals: ['genesis-orchestrator'],
          allowedEntities: ['society', 'role']
        },
        {
          type: TransactionType.LCT_CREATE,
          requiredWitnesses: 3,
          allowedEntities: ['society', 'role']
        },
        {
          type: TransactionType.ATP_TRANSFER,
          requiredWitnesses: 2,
          maxAmount: 1000
        },
        {
          type: TransactionType.ROLE_CREATE,
          requiredWitnesses: 3,
          requiredApprovals: ['genesis-orchestrator', 'queen']
        }
      ],
      
      economicRules: [
        {
          rule: 'Daily ATP allocation limit',
          parameter: 'daily_atp_max',
          value: 1000,
          enforcement: 'strict'
        },
        {
          rule: 'ADP generation ratio',
          parameter: 'adp_multiplier_max',
          value: 2,
          enforcement: 'flexible'
        },
        {
          rule: 'Minimum witness reward',
          parameter: 'witness_reward_min',
          value: 1,
          enforcement: 'strict'
        }
      ],
      
      governanceRules: [
        {
          rule: 'Law update approval',
          quorum: 66,  // 66% required
          vetoRights: ['genesis-orchestrator', 'law-oracle']
        },
        {
          rule: 'Role creation approval',
          quorum: 51,  // Simple majority
          vetoRights: ['law-oracle']
        },
        {
          rule: 'Economic rule change',
          quorum: 75,  // 75% required for economic changes
          vetoRights: ['genesis-orchestrator']
        }
      ]
    };
  }
  
  /**
   * Validate a transaction against society laws
   */
  async validateTransaction(transaction: Web4Transaction): Promise<OracleDecision> {
    const decisionId = `decision-${Date.now()}-${Math.random().toString(36).substring(7)}`;
    
    // Step 1: Check Web4 compliance
    const complianceResult = await this.config.complianceChecker.checkTransaction(transaction);
    
    if (!complianceResult.compliant) {
      const decision: OracleDecision = {
        id: decisionId,
        timestamp: Date.now(),
        transaction,
        decision: 'reject',
        reason: `Web4 compliance failed: ${complianceResult.violations[0].description}`,
        complianceResult,
        oracleSignature: this.signDecision(decisionId)
      };
      
      this.decisions.set(decisionId, decision);
      return decision;
    }
    
    // Step 2: Check transaction rules
    const rule = this.config.laws.transactionRules.find(r => r.type === transaction.type);
    
    if (rule) {
      // Check witness requirement
      if (transaction.witnesses.length < rule.requiredWitnesses) {
        const decision: OracleDecision = {
          id: decisionId,
          timestamp: Date.now(),
          transaction,
          decision: 'reject',
          reason: `Insufficient witnesses: ${transaction.witnesses.length}/${rule.requiredWitnesses}`,
          complianceResult,
          oracleSignature: this.signDecision(decisionId)
        };
        
        this.decisions.set(decisionId, decision);
        return decision;
      }
      
      // Check required approvals
      if (rule.requiredApprovals) {
        const hasApprovals = rule.requiredApprovals.every(approver =>
          transaction.witnesses.some(w => w.lct.includes(approver))
        );
        
        if (!hasApprovals) {
          const decision: OracleDecision = {
            id: decisionId,
            timestamp: Date.now(),
            transaction,
            decision: 'conditional',
            reason: 'Awaiting required approvals',
            conditions: rule.requiredApprovals.map(a => `Approval from ${a}`),
            complianceResult,
            oracleSignature: this.signDecision(decisionId)
          };
          
          this.decisions.set(decisionId, decision);
          return decision;
        }
      }
      
      // Check amount limits for ATP transfers
      if (rule.maxAmount && transaction.type === TransactionType.ATP_TRANSFER) {
        const amount = transaction.data.atp?.amount || 0;
        if (amount > rule.maxAmount) {
          const decision: OracleDecision = {
            id: decisionId,
            timestamp: Date.now(),
            transaction,
            decision: 'reject',
            reason: `ATP transfer exceeds limit: ${amount}/${rule.maxAmount}`,
            complianceResult,
            oracleSignature: this.signDecision(decisionId)
          };
          
          this.decisions.set(decisionId, decision);
          return decision;
        }
      }
    }
    
    // Step 3: Check role permissions
    const entityType = this.extractEntityType(transaction.from);
    const permissions = this.config.laws.rolePermissions.get(entityType);
    
    if (permissions) {
      const action = this.extractAction(transaction);
      const hasPermission = permissions.some(p => p.action === action);
      
      if (!hasPermission) {
        const decision: OracleDecision = {
          id: decisionId,
          timestamp: Date.now(),
          transaction,
          decision: 'reject',
          reason: `Entity ${entityType} lacks permission for ${action}`,
          complianceResult,
          oracleSignature: this.signDecision(decisionId)
        };
        
        this.decisions.set(decisionId, decision);
        return decision;
      }
    }
    
    // Step 4: Check economic rules
    if (transaction.type.startsWith('atp_') || transaction.type.startsWith('adp_')) {
      const economicViolation = this.checkEconomicRules(transaction);
      if (economicViolation) {
        const decision: OracleDecision = {
          id: decisionId,
          timestamp: Date.now(),
          transaction,
          decision: 'reject',
          reason: economicViolation,
          complianceResult,
          oracleSignature: this.signDecision(decisionId)
        };
        
        this.decisions.set(decisionId, decision);
        return decision;
      }
    }
    
    // Transaction approved!
    const decision: OracleDecision = {
      id: decisionId,
      timestamp: Date.now(),
      transaction,
      decision: 'approve',
      reason: 'Transaction complies with all society laws',
      complianceResult,
      oracleSignature: this.signDecision(decisionId)
    };
    
    this.decisions.set(decisionId, decision);
    
    // Record in blockchain
    await this.recordDecision(decision);
    
    return decision;
  }
  
  /**
   * Create a dispute
   */
  createDispute(plaintiff: string, defendant: string, claim: string, evidence: any[]): Dispute {
    const disputeId = `dispute-${Date.now()}-${Math.random().toString(36).substring(7)}`;
    
    const dispute: Dispute = {
      id: disputeId,
      plaintiff,
      defendant,
      claim,
      evidence,
      status: 'pending'
    };
    
    this.disputes.set(disputeId, dispute);
    return dispute;
  }
  
  /**
   * Resolve a dispute
   */
  async resolveDispute(disputeId: string): Promise<Dispute> {
    const dispute = this.disputes.get(disputeId);
    if (!dispute) {
      throw new Error(`Dispute ${disputeId} not found`);
    }
    
    dispute.status = 'investigating';
    
    // Check transaction history of both parties
    const plaintiffHistory = this.config.blockchain.getTransactionHistory(dispute.plaintiff);
    const defendantHistory = this.config.blockchain.getTransactionHistory(dispute.defendant);
    
    // Check compliance of both parties
    const plaintiffCompliance = await this.config.complianceChecker.createComplianceReport(
      dispute.plaintiff,
      plaintiffHistory,
      []
    );
    
    const defendantCompliance = await this.config.complianceChecker.createComplianceReport(
      dispute.defendant,
      defendantHistory,
      []
    );
    
    // Make decision based on compliance scores
    if (plaintiffCompliance.overallCompliance > defendantCompliance.overallCompliance) {
      dispute.resolution = `Ruling in favor of plaintiff. Defendant shows lower compliance (${defendantCompliance.overallCompliance}%)`;
    } else if (defendantCompliance.overallCompliance > plaintiffCompliance.overallCompliance) {
      dispute.resolution = `Ruling in favor of defendant. Plaintiff shows lower compliance (${plaintiffCompliance.overallCompliance}%)`;
    } else {
      dispute.resolution = 'No clear violation found. Parties encouraged to reach mutual agreement.';
    }
    
    dispute.status = 'resolved';
    this.disputes.set(disputeId, dispute);
    
    return dispute;
  }
  
  /**
   * Update society laws (requires governance approval)
   */
  async updateLaws(
    proposedChanges: Partial<SocietyLaws>,
    proposer: string,
    supporters: string[]
  ): Promise<boolean> {
    // Check if proposer has rights
    if (!proposer.includes('orchestrator') && !proposer.includes('oracle')) {
      console.log('Proposer lacks authority to change laws');
      return false;
    }
    
    // Check quorum
    const governanceRule = this.config.laws.governanceRules.find(r => 
      r.rule === 'Law update approval'
    );
    
    if (governanceRule) {
      const totalCitizens = 32; // In production, would query blockchain
      const supportPercentage = (supporters.length / totalCitizens) * 100;
      
      if (supportPercentage < governanceRule.quorum) {
        console.log(`Insufficient support: ${supportPercentage}% < ${governanceRule.quorum}%`);
        return false;
      }
      
      // Check for vetos
      const vetoUsed = governanceRule.vetoRights?.some(vetoEntity =>
        !supporters.some(s => s.includes(vetoEntity))
      );
      
      if (vetoUsed) {
        console.log('Law change vetoed');
        return false;
      }
    }
    
    // Apply changes
    Object.assign(this.config.laws, proposedChanges);
    this.config.laws.version = this.incrementVersion(this.config.laws.version);
    
    // Record in blockchain
    const lawUpdateTx = this.config.blockchain.createTransaction({
      type: TransactionType.SOCIETY_LAW,
      from: this.config.oracleLCT,
      to: this.config.societyLCT,
      data: {
        action: 'update_laws',
        description: `Laws updated to version ${this.config.laws.version}`,
        society: {
          societyId: this.config.societyLCT,
          action: 'update_law',
          lawUpdate: proposedChanges
        }
      },
      signature: this.signDecision('law-update'),
      witnesses: supporters.map(s => ({
        lct: s,
        signature: `support-${s}`,
        timestamp: Date.now(),
        confidence: 1.0
      })),
      validated: true
    });
    
    return true;
  }
  
  /**
   * Check economic rules
   */
  private checkEconomicRules(transaction: Web4Transaction): string | null {
    for (const rule of this.config.laws.economicRules) {
      if (rule.parameter === 'daily_atp_max' && transaction.type === TransactionType.ATP_TRANSFER) {
        // In production, would check daily totals
        const amount = transaction.data.atp?.amount || 0;
        if (amount > rule.value) {
          return `Exceeds daily ATP limit: ${amount}/${rule.value}`;
        }
      }
      
      if (rule.parameter === 'adp_multiplier_max' && transaction.type === TransactionType.ADP_GENERATE) {
        const ratio = (transaction.data.adp?.adpGenerated || 0) / 
                     (transaction.data.adp?.atpConsumed || 1);
        if (ratio > rule.value) {
          return `ADP generation ratio too high: ${ratio}/${rule.value}`;
        }
      }
    }
    
    return null;
  }
  
  /**
   * Extract entity type from LCT
   */
  private extractEntityType(lct: string): string {
    const parts = lct.split(':');
    if (parts.length >= 4) {
      // Check for specific roles
      if (parts[3].includes('orchestrator')) return 'genesis-orchestrator';
      if (parts[3].includes('queen')) return 'queen';
      if (parts[3].includes('worker')) return 'worker';
      if (parts[3].includes('oracle')) return 'law-oracle';
    }
    return parts[2] || 'unknown';
  }
  
  /**
   * Extract action from transaction
   */
  private extractAction(transaction: Web4Transaction): string {
    const typeMap: Record<string, string> = {
      'lct_create': 'create',
      'role_create': 'spawn',
      'atp_transfer': 'allocate',
      'role_execute': 'execute',
      'witness_action': 'witness'
    };
    
    return typeMap[transaction.type] || 'unknown';
  }
  
  /**
   * Sign a decision
   */
  private signDecision(decisionId: string): string {
    // In production, would use Ed25519 signature
    return `oracle-sig-${decisionId.substring(0, 8)}`;
  }
  
  /**
   * Record decision in blockchain
   */
  private async recordDecision(decision: OracleDecision): Promise<void> {
    this.config.blockchain.createTransaction({
      type: TransactionType.WITNESS_VALIDATE,
      from: this.config.oracleLCT,
      to: decision.transaction?.from,
      data: {
        action: 'validate_transaction',
        description: `Oracle decision: ${decision.decision}`,
        witness: {
          action: decision.reason,
          rolePerformed: 'law-oracle',
          timestamp: decision.timestamp,
          evidence: decision
        }
      },
      signature: decision.oracleSignature,
      witnesses: [],
      validated: true
    });
  }
  
  /**
   * Increment version number
   */
  private incrementVersion(version: string): string {
    const parts = version.split('.');
    parts[2] = (parseInt(parts[2]) + 1).toString();
    return parts.join('.');
  }
  
  /**
   * Get oracle statistics
   */
  getStatistics(): {
    totalDecisions: number;
    approvedTransactions: number;
    rejectedTransactions: number;
    pendingDisputes: number;
    resolvedDisputes: number;
    complianceRate: number;
  } {
    const decisions = Array.from(this.decisions.values());
    const disputes = Array.from(this.disputes.values());
    
    const approved = decisions.filter(d => d.decision === 'approve').length;
    const rejected = decisions.filter(d => d.decision === 'reject').length;
    
    return {
      totalDecisions: decisions.length,
      approvedTransactions: approved,
      rejectedTransactions: rejected,
      pendingDisputes: disputes.filter(d => d.status === 'pending').length,
      resolvedDisputes: disputes.filter(d => d.status === 'resolved').length,
      complianceRate: decisions.length > 0 ? (approved / decisions.length) * 100 : 100
    };
  }
}

// Export for use in ACT
export default LawOracle;