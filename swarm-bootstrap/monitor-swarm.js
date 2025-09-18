#!/usr/bin/env node

/**
 * ACT Swarm Progress Monitor
 * Real-time visibility into the fractal swarm building ACT
 */

const fs = require('fs');
const path = require('path');
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

// Colors for terminal output
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  dim: '\x1b[2m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
  white: '\x1b[37m'
};

// Load swarm configuration
const SWARM_CONFIG = require('./swarm-config.json');

// Memory paths
const MEMORY_BASE = './swarm-memory';
const WITNESS_LOG = path.join(MEMORY_BASE, 'witness/activity.jsonl');
const ATP_LEDGER = path.join(MEMORY_BASE, 'economy/atp-ledger.json');
const PROGRESS_FILE = path.join(MEMORY_BASE, 'implementation/progress.json');

/**
 * Display swarm hierarchy with current status
 */
async function displaySwarmHierarchy() {
  console.log(`\n${colors.bright}${colors.cyan}🌟 ACT Development Swarm Status${colors.reset}`);
  console.log('=' + '='.repeat(60));
  
  // Genesis Orchestrator
  console.log(`\n${colors.bright}👑 Genesis Orchestrator${colors.reset}`);
  console.log(`   Name: ${SWARM_CONFIG.genesis.name}`);
  console.log(`   Mode: Delegative`);
  console.log(`   Daily ATP: ${SWARM_CONFIG.genesis.budget.daily}`);
  console.log(`   Status: ${colors.green}● Active${colors.reset}`);
  
  // Domain Queens
  console.log(`\n${colors.bright}🏗️ Domain Queens (7)${colors.reset}`);
  for (const queen of SWARM_CONFIG.queens) {
    const status = await getQueenStatus(queen.name);
    const statusColor = status === 'active' ? colors.green : colors.yellow;
    console.log(`   ${statusColor}●${colors.reset} ${queen.name}`);
    console.log(`      Domain: ${queen.domain}`);
    console.log(`      Workers: ${queen.workers.length}`);
    console.log(`      ATP Budget: ${queen.budget}`);
  }
  
  // Worker summary
  const totalWorkers = SWARM_CONFIG.queens.reduce((sum, q) => sum + q.workers.length, 0);
  console.log(`\n${colors.bright}🔧 Worker Roles${colors.reset}`);
  console.log(`   Total: ${totalWorkers} workers`);
  console.log(`   Mode: Responsive`);
  console.log(`   Status: Awaiting tasks`);
}

/**
 * Display ATP economy status
 */
async function displayATPEconomy() {
  console.log(`\n${colors.bright}${colors.yellow}💰 ATP Economy${colors.reset}`);
  console.log('=' + '='.repeat(60));
  
  try {
    let atpData = { treasury: 10000, spent: 0, allocated: 0 };
    if (fs.existsSync(ATP_LEDGER)) {
      atpData = JSON.parse(fs.readFileSync(ATP_LEDGER, 'utf8'));
    }
    
    const remaining = atpData.treasury - atpData.spent - atpData.allocated;
    const efficiency = atpData.spent > 0 ? 
      ((atpData.adp_generated || 0) / atpData.spent * 100).toFixed(1) : '0.0';
    
    console.log(`   Treasury:  ${atpData.treasury} ATP`);
    console.log(`   Allocated: ${atpData.allocated} ATP`);
    console.log(`   Spent:     ${atpData.spent} ATP`);
    console.log(`   Remaining: ${remaining} ATP`);
    console.log(`   ADP Efficiency: ${efficiency}%`);
    
    // Show recent transactions
    console.log(`\n   Recent ATP Transactions:`);
    const transactions = atpData.recent_transactions || [];
    transactions.slice(-5).forEach(tx => {
      const icon = tx.type === 'spend' ? '↓' : '↑';
      console.log(`     ${icon} ${tx.amount} ATP - ${tx.description}`);
    });
  } catch (error) {
    console.log(`   ${colors.dim}No ATP data available yet${colors.reset}`);
  }
}

/**
 * Display implementation progress
 */
async function displayProgress() {
  console.log(`\n${colors.bright}${colors.green}📊 Implementation Progress${colors.reset}`);
  console.log('=' + '='.repeat(60));
  
  try {
    let progress = {
      phases: {
        foundation: { complete: 0, total: 10 },
        protocol: { complete: 0, total: 8 },
        integration: { complete: 0, total: 12 },
        interface: { complete: 0, total: 6 }
      }
    };
    
    if (fs.existsSync(PROGRESS_FILE)) {
      progress = JSON.parse(fs.readFileSync(PROGRESS_FILE, 'utf8'));
    }
    
    // Display progress bars for each phase
    for (const [phase, data] of Object.entries(progress.phases)) {
      const percent = (data.complete / data.total * 100).toFixed(0);
      const filled = Math.floor(percent / 5);
      const empty = 20 - filled;
      const bar = '█'.repeat(filled) + '░'.repeat(empty);
      
      console.log(`   ${phase.padEnd(12)} [${bar}] ${percent}% (${data.complete}/${data.total})`);
    }
    
    // Overall progress
    const totalComplete = Object.values(progress.phases)
      .reduce((sum, p) => sum + p.complete, 0);
    const totalTasks = Object.values(progress.phases)
      .reduce((sum, p) => sum + p.total, 0);
    const overallPercent = (totalComplete / totalTasks * 100).toFixed(0);
    
    console.log(`\n   ${colors.bright}Overall Progress: ${overallPercent}%${colors.reset}`);
    console.log(`   Tasks Completed: ${totalComplete}/${totalTasks}`);
    
    // Estimated completion
    const startTime = progress.start_time || Date.now();
    const elapsed = Date.now() - startTime;
    const rate = totalComplete > 0 ? totalComplete / (elapsed / (1000 * 60 * 60 * 24)) : 0;
    const remaining = totalTasks - totalComplete;
    const daysRemaining = rate > 0 ? Math.ceil(remaining / rate) : 28;
    
    console.log(`   Estimated Completion: ${daysRemaining} days`);
    
  } catch (error) {
    console.log(`   ${colors.dim}No progress data available yet${colors.reset}`);
  }
}

/**
 * Display recent witness activity
 */
async function displayWitnessActivity() {
  console.log(`\n${colors.bright}${colors.magenta}👁️ Witness Activity${colors.reset}`);
  console.log('=' + '='.repeat(60));
  
  try {
    if (fs.existsSync(WITNESS_LOG)) {
      const lines = fs.readFileSync(WITNESS_LOG, 'utf8').trim().split('\n');
      const recentActivity = lines.slice(-10).map(line => JSON.parse(line));
      
      console.log(`   Recent Actions (Last 10):`);
      recentActivity.forEach(action => {
        const time = new Date(action.timestamp).toLocaleTimeString();
        const roleShort = action.role_lct?.split(':').pop() || 'unknown';
        console.log(`   ${time} | ${roleShort.padEnd(20)} | ${action.action}`);
      });
      
      // Activity stats
      const allActivity = lines.map(line => JSON.parse(line));
      const roleStats = {};
      allActivity.forEach(action => {
        const role = action.role_lct?.split(':').pop() || 'unknown';
        roleStats[role] = (roleStats[role] || 0) + 1;
      });
      
      console.log(`\n   Activity by Role:`);
      Object.entries(roleStats)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 5)
        .forEach(([role, count]) => {
          console.log(`     ${role.padEnd(25)} ${count} actions`);
        });
    } else {
      console.log(`   ${colors.dim}No witness activity yet${colors.reset}`);
    }
  } catch (error) {
    console.log(`   ${colors.dim}Error reading witness log${colors.reset}`);
  }
}

/**
 * Display current tasks
 */
async function displayCurrentTasks() {
  console.log(`\n${colors.bright}${colors.blue}🎯 Current Tasks${colors.reset}`);
  console.log('=' + '='.repeat(60));
  
  // Simulated current tasks (would come from swarm memory)
  const currentTasks = [
    { queen: 'LCT-Infrastructure-Queen', task: 'Implementing Ed25519 binding', status: 'in_progress', workers: 3 },
    { queen: 'ACP-Protocol-Queen', task: 'Designing agent plan structure', status: 'planning', workers: 2 },
    { queen: 'Demo-Society-Queen', task: 'Setting up law oracle', status: 'queued', workers: 0 },
    { queen: 'MCP-Bridge-Queen', task: 'Claude MCP integration', status: 'queued', workers: 0 },
    { queen: 'Client-Interface-Queen', task: 'Dashboard wireframes', status: 'planning', workers: 1 },
    { queen: 'ATP-Economy-Queen', task: 'Token mechanics design', status: 'in_progress', workers: 2 }
  ];
  
  currentTasks.forEach(task => {
    const statusIcon = {
      'in_progress': `${colors.green}◉${colors.reset}`,
      'planning': `${colors.yellow}◎${colors.reset}`,
      'queued': `${colors.dim}○${colors.reset}`
    }[task.status];
    
    console.log(`   ${statusIcon} ${task.queen}`);
    console.log(`      Task: ${task.task}`);
    console.log(`      Workers: ${task.workers} active`);
  });
}

/**
 * Get queen status (simulated for now)
 */
async function getQueenStatus(queenName) {
  // In production, this would check actual Claude-Flow status
  return 'active';
}

/**
 * Live monitoring mode
 */
async function liveMonitor() {
  console.clear();
  
  while (true) {
    console.log('\x1Bc'); // Clear screen
    console.log(`${colors.bright}ACT Swarm Monitor - ${new Date().toLocaleString()}${colors.reset}`);
    
    await displaySwarmHierarchy();
    await displayProgress();
    await displayCurrentTasks();
    await displayATPEconomy();
    await displayWitnessActivity();
    
    console.log(`\n${colors.dim}Refreshing every 5 seconds... (Ctrl+C to exit)${colors.reset}`);
    
    await new Promise(resolve => setTimeout(resolve, 5000));
  }
}

/**
 * Main execution
 */
async function main() {
  const args = process.argv.slice(2);
  const mode = args[0] || 'snapshot';
  
  if (mode === 'live') {
    await liveMonitor();
  } else {
    // Single snapshot
    await displaySwarmHierarchy();
    await displayProgress();
    await displayCurrentTasks();
    await displayATPEconomy();
    await displayWitnessActivity();
    
    console.log(`\n${colors.dim}Run 'node monitor-swarm.js live' for real-time monitoring${colors.reset}\n`);
  }
}

// Error handling
process.on('SIGINT', () => {
  console.log(`\n${colors.yellow}Monitoring stopped${colors.reset}`);
  process.exit(0);
});

// Run if executed directly
if (require.main === module) {
  main().catch(error => {
    console.error(`${colors.red}Error: ${error.message}${colors.reset}`);
    process.exit(1);
  });
}

module.exports = {
  displaySwarmHierarchy,
  displayProgress,
  displayATPEconomy,
  displayWitnessActivity
};