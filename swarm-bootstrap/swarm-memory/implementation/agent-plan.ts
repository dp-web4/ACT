// ACP Agent Plan Structure
export interface AgentPlan {
  id: string;
  owner_lct: string;
  triggers: Trigger[];
  intents: Intent[];
  decisions: Decision[];
  created_at: number;
}

export interface Trigger {
  type: 'event' | 'time' | 'condition';
  specification: any;
}