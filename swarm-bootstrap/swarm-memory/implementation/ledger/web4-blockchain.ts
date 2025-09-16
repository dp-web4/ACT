/**
 * Web4-Compliant Immutable Blockchain Ledger
 * Core ledger for ACT Society
 * 
 * This implements witnessed presence protocol with:
 * - Immutable blocks with cryptographic linking
 * - LCT-based witness signatures
 * - ATP/ADP transaction records
 * - MRH (Markov Relevancy Horizon) tracking
 */

import { createHash } from 'crypto';
import { ed25519 } from '@noble/curves/ed25519';

// Web4 Block Structure
export interface Web4Block {
  // Block metadata
  index: number;
  timestamp: number;
  previousHash: string;
  hash: string;
  nonce: number;
  
  // Web4 specific fields
  society: string;                    // Society LCT that owns this ledger
  lawOracle: string;                   // Law oracle LCT for validation
  witnessedBy: WitnessSignature[];    // Multiple witnesses per block
  
  // Transactions in this block
  transactions: Web4Transaction[];
  
  // MRH tracking
  mrh: {
    bound: string[];      // Entities bound to this block
    paired: string[];     // Entities paired during block creation
    witnessing: string[]; // Entities witnessing this block
    broadcast: string[];  // Entities that received broadcast
  };
}

// Web4 Transaction Types
export interface Web4Transaction {
  id: string;
  type: TransactionType;
  timestamp: number;
  
  // Core Web4 fields
  from: string;         // LCT of sender
  to?: string;          // LCT of receiver (optional)
  
  // Transaction data
  data: {
    // Common fields
    action: string;
    description: string;
    
    // Type-specific fields
    lct?: LCTTransaction;
    atp?: ATPTransaction;
    adp?: ADPProof;
    witness?: WitnessRecord;
    role?: RoleTransaction;
    society?: SocietyTransaction;
  };
  
  // Witness and validation
  signature: string;              // Ed25519 signature from sender
  witnesses: WitnessSignature[];  // Required witnesses
  validated: boolean;             // Law oracle validation
}

export enum TransactionType {
  // LCT operations
  LCT_CREATE = 'lct_create',
  LCT_BIND = 'lct_bind',
  LCT_PAIR = 'lct_pair',
  LCT_WITNESS = 'lct_witness',
  
  // ATP/ADP operations
  ATP_TRANSFER = 'atp_transfer',
  ATP_ALLOCATE = 'atp_allocate',
  ADP_GENERATE = 'adp_generate',
  ADP_CLAIM = 'adp_claim',
  
  // Role operations
  ROLE_CREATE = 'role_create',
  ROLE_ASSIGN = 'role_assign',
  ROLE_EXECUTE = 'role_execute',
  ROLE_COMPLETE = 'role_complete',
  
  // Society operations
  SOCIETY_CREATE = 'society_create',
  SOCIETY_JOIN = 'society_join',
  SOCIETY_LEAVE = 'society_leave',
  SOCIETY_LAW = 'society_law',
  
  // Witness operations
  WITNESS_ACTION = 'witness_action',
  WITNESS_VALIDATE = 'witness_validate'
}

// Specific transaction types
export interface LCTTransaction {
  lctId: string;
  entityType: 'human' | 'ai' | 'role' | 'society' | 'dictionary';
  publicKey: string;
  birthCertificate?: {
    society: string;
    rights: string[];
    responsibilities: string[];
    initialATP: number;
  };
}

export interface ATPTransaction {
  amount: number;
  purpose: string;
  maxFee?: number;
  priority?: 'low' | 'normal' | 'high' | 'critical';
}

export interface ADPProof {
  taskCompleted: string;
  atpConsumed: number;
  adpGenerated: number;
  r6Compliance: {
    rules: boolean;
    roles: boolean;
    request: boolean;
    reference: boolean;
    resource: boolean;
    result: boolean;
  };
}

export interface WitnessRecord {
  action: string;
  rolePerformed: string;
  timestamp: number;
  evidence?: any;
}

export interface WitnessSignature {
  lct: string;           // LCT of witness
  signature: string;     // Ed25519 signature
  timestamp: number;     // When witnessed
  confidence: number;    // 0-1 confidence score
}

export interface RoleTransaction {
  roleId: string;
  action: 'create' | 'assign' | 'execute' | 'complete';
  r6Rules?: any;
  performer?: string;  // LCT performing role
  task?: string;
  result?: any;
}

export interface SocietyTransaction {
  societyId: string;
  action: 'create' | 'join' | 'leave' | 'update_law';
  citizenRole?: string;
  lawUpdate?: any;
}

/**
 * Web4-Compliant Blockchain Implementation
 */
export class Web4Blockchain {
  private chain: Web4Block[] = [];
  private pendingTransactions: Web4Transaction[] = [];
  private miningReward = 100; // ATP reward for mining
  private difficulty = 2;     // Proof of work difficulty
  
  // Web4 specific
  private society: string;
  private lawOracle: string;
  private witnesses: Set<string> = new Set();
  
  constructor(society: string, lawOracle: string) {
    this.society = society;
    this.lawOracle = lawOracle;
    this.chain = [this.createGenesisBlock()];
  }
  
  /**
   * Create the genesis block for ACT Society
   */
  private createGenesisBlock(): Web4Block {
    const genesisTransaction: Web4Transaction = {
      id: 'genesis-tx',
      type: TransactionType.SOCIETY_CREATE,
      timestamp: Date.now(),
      from: 'lct:web4:genesis',
      data: {
        action: 'create_society',
        description: 'ACT Development Collective Genesis',
        society: {
          societyId: this.society,
          action: 'create',
          lawUpdate: {
            name: 'ACT Development Collective',
            version: '1.0.0',
            rules: ['Web4 compliance', 'Trust-native', 'Witnessed presence']
          }
        }
      },
      signature: 'genesis-signature',
      witnesses: [],
      validated: true
    };
    
    return {
      index: 0,
      timestamp: Date.now(),
      previousHash: '0',
      hash: this.calculateHash({
        index: 0,
        timestamp: Date.now(),
        previousHash: '0',
        transactions: [genesisTransaction],
        nonce: 0
      } as any),
      nonce: 0,
      society: this.society,
      lawOracle: this.lawOracle,
      witnessedBy: [{
        lct: 'lct:web4:genesis',
        signature: 'genesis-witness',
        timestamp: Date.now(),
        confidence: 1.0
      }],
      transactions: [genesisTransaction],
      mrh: {
        bound: [this.society],
        paired: [],
        witnessing: ['lct:web4:genesis'],
        broadcast: []
      }
    };
  }
  
  /**
   * Get the latest block
   */
  getLatestBlock(): Web4Block {
    return this.chain[this.chain.length - 1];
  }
  
  /**
   * Calculate hash for a block
   */
  calculateHash(block: Partial<Web4Block>): string {
    const data = JSON.stringify({
      index: block.index,
      timestamp: block.timestamp,
      previousHash: block.previousHash,
      transactions: block.transactions,
      nonce: block.nonce,
      society: block.society
    });
    
    return createHash('sha256').update(data).digest('hex');
  }
  
  /**
   * Mine a new block (Proof of Work)
   */
  mineBlock(block: Web4Block): void {
    while (block.hash.substring(0, this.difficulty) !== Array(this.difficulty + 1).join('0')) {
      block.nonce++;
      block.hash = this.calculateHash(block);
    }
    
    console.log(`Block mined: ${block.hash}`);
  }
  
  /**
   * Create a new Web4 transaction
   */
  createTransaction(transaction: Omit<Web4Transaction, 'id' | 'timestamp'>): Web4Transaction {
    const tx: Web4Transaction = {
      ...transaction,
      id: `tx-${Date.now()}-${Math.random().toString(36).substring(7)}`,
      timestamp: Date.now()
    };
    
    // Validate with law oracle
    tx.validated = this.validateWithLawOracle(tx);
    
    this.pendingTransactions.push(tx);
    return tx;
  }
  
  /**
   * Mine pending transactions into a new block
   */
  minePendingTransactions(minerLCT: string, witnesses: WitnessSignature[]): Web4Block | null {
    if (this.pendingTransactions.length === 0) {
      return null;
    }
    
    const block: Web4Block = {
      index: this.chain.length,
      timestamp: Date.now(),
      previousHash: this.getLatestBlock().hash,
      hash: '',
      nonce: 0,
      society: this.society,
      lawOracle: this.lawOracle,
      witnessedBy: witnesses,
      transactions: [...this.pendingTransactions],
      mrh: {
        bound: [this.society, minerLCT],
        paired: this.extractPairedEntities(this.pendingTransactions),
        witnessing: witnesses.map(w => w.lct),
        broadcast: [] // Will be filled when broadcast
      }
    };
    
    // Calculate initial hash
    block.hash = this.calculateHash(block);
    
    // Mine the block
    this.mineBlock(block);
    
    // Add to chain
    this.chain.push(block);
    
    // Clear pending transactions
    this.pendingTransactions = [];
    
    // Reward miner with ATP
    this.createTransaction({
      type: TransactionType.ATP_TRANSFER,
      from: this.society,
      to: minerLCT,
      data: {
        action: 'mining_reward',
        description: `Mining reward for block ${block.index}`,
        atp: {
          amount: this.miningReward,
          purpose: 'Block mining reward',
          priority: 'normal'
        }
      },
      signature: 'system-signature',
      witnesses: [],
      validated: true
    });
    
    return block;
  }
  
  /**
   * Validate transaction with law oracle
   */
  private validateWithLawOracle(transaction: Web4Transaction): boolean {
    // Simulate law oracle validation
    // In production, this would call the actual law oracle role
    
    // Check transaction type permissions
    switch (transaction.type) {
      case TransactionType.SOCIETY_CREATE:
        // Only genesis can create societies initially
        return transaction.from === 'lct:web4:genesis' || 
               transaction.from === this.society;
      
      case TransactionType.LCT_CREATE:
        // Must be from a society or authorized role
        return transaction.from.includes('society') || 
               transaction.from.includes('role');
      
      case TransactionType.ATP_TRANSFER:
        // Must have sufficient ATP (would check balance)
        return true;
      
      case TransactionType.WITNESS_ACTION:
        // Anyone can witness
        return true;
      
      default:
        // Default validation
        return true;
    }
  }
  
  /**
   * Extract paired entities from transactions
   */
  private extractPairedEntities(transactions: Web4Transaction[]): string[] {
    const paired: string[] = [];
    
    for (const tx of transactions) {
      if (tx.type === TransactionType.LCT_PAIR && tx.to) {
        paired.push(tx.from, tx.to);
      }
    }
    
    return [...new Set(paired)];
  }
  
  /**
   * Validate the entire blockchain
   */
  isChainValid(): boolean {
    for (let i = 1; i < this.chain.length; i++) {
      const currentBlock = this.chain[i];
      const previousBlock = this.chain[i - 1];
      
      // Validate current block hash
      if (currentBlock.hash !== this.calculateHash(currentBlock)) {
        console.error(`Invalid hash at block ${i}`);
        return false;
      }
      
      // Validate link to previous block
      if (currentBlock.previousHash !== previousBlock.hash) {
        console.error(`Invalid previous hash at block ${i}`);
        return false;
      }
      
      // Validate proof of work
      if (currentBlock.hash.substring(0, this.difficulty) !== 
          Array(this.difficulty + 1).join('0')) {
        console.error(`Invalid proof of work at block ${i}`);
        return false;
      }
      
      // Validate witnesses (must have at least one)
      if (currentBlock.witnessedBy.length === 0) {
        console.error(`No witnesses for block ${i}`);
        return false;
      }
      
      // Validate all transactions
      for (const tx of currentBlock.transactions) {
        if (!tx.validated) {
          console.error(`Invalid transaction ${tx.id} in block ${i}`);
          return false;
        }
      }
    }
    
    return true;
  }
  
  /**
   * Get balance for an LCT
   */
  getBalance(lct: string): { atp: number; adp: number } {
    let atpBalance = 0;
    let adpBalance = 0;
    
    for (const block of this.chain) {
      for (const tx of block.transactions) {
        // ATP transactions
        if (tx.type === TransactionType.ATP_TRANSFER && tx.data.atp) {
          if (tx.from === lct) {
            atpBalance -= tx.data.atp.amount;
          }
          if (tx.to === lct) {
            atpBalance += tx.data.atp.amount;
          }
        }
        
        // ADP generation
        if (tx.type === TransactionType.ADP_GENERATE && tx.to === lct && tx.data.adp) {
          adpBalance += tx.data.adp.adpGenerated;
        }
      }
    }
    
    return { atp: atpBalance, adp: adpBalance };
  }
  
  /**
   * Get transaction history for an LCT
   */
  getTransactionHistory(lct: string): Web4Transaction[] {
    const transactions: Web4Transaction[] = [];
    
    for (const block of this.chain) {
      for (const tx of block.transactions) {
        if (tx.from === lct || tx.to === lct) {
          transactions.push(tx);
        }
      }
    }
    
    return transactions;
  }
  
  /**
   * Export chain to JSON
   */
  exportChain(): string {
    return JSON.stringify(this.chain, null, 2);
  }
  
  /**
   * Import chain from JSON
   */
  importChain(chainData: string): boolean {
    try {
      const importedChain = JSON.parse(chainData);
      
      // Validate imported chain
      const tempBlockchain = new Web4Blockchain(this.society, this.lawOracle);
      tempBlockchain.chain = importedChain;
      
      if (tempBlockchain.isChainValid()) {
        this.chain = importedChain;
        return true;
      }
      
      return false;
    } catch (error) {
      console.error('Failed to import chain:', error);
      return false;
    }
  }
}

// Export for use in ACT
export default Web4Blockchain;