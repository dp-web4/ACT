// LCT Token Implementation
export interface LCT {
  id: string;
  entity_type: 'human' | 'ai' | 'role' | 'society' | 'dictionary';
  public_key: string;
  binding_signature: string;
  mrh: {
    bound: string[];
    paired: string[];
    witnessing: string[];
  };
  created_at: number;
}