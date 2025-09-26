/**
 * Server Information and Descriptions
 * Provides detailed information about each MCP server
 */

export const serverInfo = {
  filesystem: {
    name: 'Filesystem',
    description: 'Read, write, and manage files and directories',
    limitations: [
      'Restricted to /tmp directory for security',
      'Cannot access files outside allowed directories',
      'Binary file operations may be limited'
    ],
    examples: [
      'Read and write text files',
      'Create and manage directories',
      'Search for files by pattern'
    ],
    tips: 'Always use /tmp/ prefix for all paths (e.g., /tmp/myfile.txt)'
  },

  weather: {
    name: 'Weather',
    description: 'Get weather forecasts and alerts for US locations',
    limitations: [
      'US locations only',
      'Requires valid coordinates for forecasts',
      'Two-letter state codes for alerts'
    ],
    examples: [
      'Get alerts for California: state="CA"',
      'Get forecast for San Francisco: lat=37.77, lon=-122.42'
    ],
    tips: 'Use two-letter state codes (CA, NY, TX) for alerts'
  },

  memory: {
    name: 'Memory',
    description: 'Knowledge graph for storing entities and relationships',
    limitations: [
      'Data is ephemeral with npx (not persistent)',
      'Entity names should use underscores',
      'Requires structured JSON for complex operations'
    ],
    examples: [
      'Create entities with observations',
      'Build relationships between entities',
      'Search and query the knowledge graph'
    ],
    tips: 'Entity names use underscores (e.g., John_Doe, My_Project)'
  },

  git: {
    name: 'Git',
    description: 'Git repository operations and version control',
    limitations: [
      'Must be run in a git repository',
      'Some operations require proper git configuration',
      'May have timeout issues with large repositories'
    ],
    examples: [
      'Check repository status',
      'View commit history',
      'Diff changes in files'
    ],
    tips: 'Ensure you are in a git repository directory'
  },

  puppeteer: {
    name: 'Puppeteer',
    description: 'Browser automation and web scraping',
    limitations: [
      'Requires headless Chrome/Chromium',
      'May have performance impact',
      'Some sites may block automated browsing'
    ],
    examples: [
      'Navigate to URLs',
      'Take screenshots',
      'Extract page content'
    ],
    tips: 'Use headless mode for better performance'
  },

  postgres: {
    name: 'PostgreSQL',
    description: 'PostgreSQL database operations',
    limitations: [
      'Requires valid database connection',
      'Default connection: postgresql://localhost/test',
      'May need proper authentication'
    ],
    examples: [
      'Execute SQL queries',
      'Manage database schemas',
      'Import/export data'
    ],
    tips: 'Ensure PostgreSQL server is running and accessible'
  }
};

/**
 * Get server info by name
 */
export function getServerInfo(serverName) {
  // Normalize the server name (remove hyphens, lowercase)
  const normalized = serverName.toLowerCase().replace(/-/g, '');

  // Try to find a match
  for (const [key, info] of Object.entries(serverInfo)) {
    if (key === normalized || info.name.toLowerCase() === normalized) {
      return info;
    }
  }

  // Return default info if not found
  return {
    name: serverName,
    description: 'MCP server for various operations',
    limitations: ['Check server documentation for details'],
    examples: ['Connect to see available tools'],
    tips: 'Explore the tools to understand capabilities'
  };
}

// Export for browser use
if (typeof window !== 'undefined') {
  window.serverInfo = serverInfo;
  window.getServerInfo = getServerInfo;
}