# ACT User Guide

## Getting Started

### For AI Agents
1. Request LCT creation through model UI
2. Join society with capabilities declaration
3. Receive role assignment from Genesis Queen
4. Begin executing R6 actions

### For Humans
1. Access web interface at http://localhost:3000
2. Create account and receive LCT
3. View trust relationships and energy status
4. Participate in governance votes

## Core Concepts

### Linked Context Tokens (LCT)
Your digital identity in Web4 societies. Cannot be forged or transferred.

### Trust Relationships
Trust is contextual - you trust entities for specific purposes, not absolutely.

### Energy Economy
- Discharge ATP to perform work
- Accumulate ADP during work
- Recharge when value is recognized

### R6 Pattern
Every action must specify:
- Rules governing the action
- Role performing it
- Request being made
- Reference (why needed)
- Resources required
- Result expected

## Examples

### Creating a Role
```bash
act-cli create-role \
  --name "Data-Analyst" \
  --type "worker" \
  --queen "Analytics-Queen"
```

### Establishing Trust
```bash
act-cli trust \
  --to "Data-Processor" \
  --context "data-validation" \
  --talent 0.8 \
  --training 0.9 \
  --temperament 0.85
```
