/**
 * Tool helpers - Provides examples and input templates for MCP tools
 */

export const toolExamples = {
  // Memory server tools
  'create_entities': {
    description: 'Create entities in the knowledge graph',
    example: {
      entities: [
        {
          name: "Dennis_Palatov",
          entityType: "person",
          observations: [
            "Creator of Web4",
            "Goes by 'dp'",
            "Working on distributed intelligence systems"
          ]
        },
        {
          name: "ACT_Tool",
          entityType: "software",
          observations: [
            "Direct MCP server interface",
            "Created September 2025",
            "Enables human interaction with MCP servers"
          ]
        }
      ]
    },
    tips: [
      'Entity names should use underscores for spaces',
      'Entity types: person, organization, software, project, concept, event',
      'Observations should be atomic facts (one fact per observation)'
    ]
  },

  'create_relations': {
    description: 'Create relationships between entities',
    example: {
      relations: [
        {
          from: "Dennis_Palatov",
          to: "ACT_Tool",
          relationType: "created"
        },
        {
          from: "ACT_Tool",
          to: "Web4",
          relationType: "implements"
        }
      ]
    },
    tips: [
      'Relations are directional (from → to)',
      'Use active voice for relation types (created, owns, manages, works_at)',
      'Both entities must exist before creating relations'
    ]
  },

  'add_observations': {
    description: 'Add new observations to existing entities',
    example: {
      observations: [
        {
          entityName: "ACT_Tool",
          contents: [
            "Supports filesystem, weather, and memory MCP servers",
            "Has three-panel UI layout"
          ]
        }
      ]
    },
    tips: [
      'Entity must already exist',
      'Each observation should be a single fact',
      'Use present tense for current facts'
    ]
  },

  'search_nodes': {
    description: 'Search for entities in the knowledge graph',
    example: {
      query: "Web4"
    },
    tips: [
      'Searches entity names, types, and observations',
      'Partial matches are supported',
      'Case-sensitive search'
    ]
  },

  'open_nodes': {
    description: 'Get specific entities by exact name',
    example: {
      names: ["Dennis_Palatov", "ACT_Tool"]
    },
    tips: [
      'Use exact entity names',
      'Returns entities and their relations',
      'Names are case-sensitive'
    ]
  },

  'read_graph': {
    description: 'Read the entire knowledge graph',
    example: {},
    tips: [
      'No parameters needed',
      'Returns all entities and relations',
      'Can be large if graph has many nodes'
    ]
  },

  'delete_entities': {
    description: 'Remove entities from the graph',
    example: {
      entityNames: ["Old_Project"]
    },
    tips: [
      'Deletes entity and all its relations',
      'Cannot be undone',
      'Use exact entity names'
    ]
  },

  'delete_observations': {
    description: 'Remove specific observations',
    example: {
      deletions: [
        {
          entityName: "ACT_Tool",
          observations: ["Outdated information"]
        }
      ]
    },
    tips: [
      'Must match observation text exactly',
      'Entity is preserved, only observations removed'
    ]
  },

  'delete_relations': {
    description: 'Remove specific relations',
    example: {
      relations: [
        {
          from: "Person_A",
          to: "Company_B",
          relationType: "works_at"
        }
      ]
    },
    tips: [
      'Must match all three fields exactly',
      'Entities are preserved, only relation removed'
    ]
  },

  // Weather server tools
  'get-alerts': {
    description: 'Get weather alerts for a US state',
    example: {
      state: "CA"
    },
    tips: [
      'Use two-letter state codes (CA, NY, TX, etc.)',
      'Returns active weather alerts',
      'US states only'
    ]
  },

  'get-forecast': {
    description: 'Get weather forecast for coordinates',
    example: {
      latitude: 37.7749,
      longitude: -122.4194,
      days: 3
    },
    tips: [
      'Latitude: -90 to 90',
      'Longitude: -180 to 180',
      'Days: 1-10 (optional, default 3)',
      'Example coords: San Francisco (37.77, -122.42)'
    ]
  },

  // Filesystem tools
  'read_file': {
    description: 'Read contents of a file',
    example: {
      path: "README.md"
    },
    tips: [
      'Path relative to server root (/tmp)',
      'Text files only',
      'Binary files will error'
    ]
  },

  'write_file': {
    description: 'Write content to a file',
    example: {
      path: "test.txt",
      content: "Hello from ACT Tool!"
    },
    tips: [
      'Path relative to server root (/tmp)',
      'Creates directories if needed',
      'Overwrites existing files'
    ]
  },

  'list_directory': {
    description: 'List files in a directory',
    example: {
      path: "/"
    },
    tips: [
      'Path relative to server root (/tmp)',
      'Returns files and subdirectories',
      'Use "/" for root directory'
    ]
  },

  'create_directory': {
    description: 'Create a new directory',
    example: {
      path: "my_folder"
    },
    tips: [
      'Path relative to server root (/tmp)',
      'Creates parent directories if needed',
      'No error if already exists'
    ]
  },

  'delete_file': {
    description: 'Delete a file or directory',
    example: {
      path: "old_file.txt"
    },
    tips: [
      'Path relative to server root (/tmp)',
      'Deletes recursively for directories',
      'Cannot be undone'
    ]
  },

  'move_file': {
    description: 'Move or rename a file',
    example: {
      source: "old_name.txt",
      destination: "new_name.txt"
    },
    tips: [
      'Paths relative to server root (/tmp)',
      'Can move between directories',
      'Overwrites destination if exists'
    ]
  },

  'search_files': {
    description: 'Search for files by pattern',
    example: {
      path: "/",
      pattern: "*.txt"
    },
    tips: [
      'Supports wildcards (* and ?)',
      'Case-sensitive by default',
      'Searches recursively'
    ]
  },

  'get_file_info': {
    description: 'Get metadata about a file',
    example: {
      path: "README.md"
    },
    tips: [
      'Returns size, permissions, timestamps',
      'Works for files and directories'
    ]
  },

  // Git tools (common ones)
  'git_status': {
    description: 'Get repository status',
    example: {},
    tips: [
      'No parameters needed',
      'Shows modified, staged, untracked files',
      'Repository must be initialized'
    ]
  },

  'git_log': {
    description: 'Get commit history',
    example: {
      max_count: 10
    },
    tips: [
      'max_count: number of commits to show',
      'Returns commit hash, author, date, message'
    ]
  },

  'git_diff': {
    description: 'Show changes in files',
    example: {
      path: "src/index.js"
    },
    tips: [
      'Leave path empty for all changes',
      'Shows unstaged changes by default'
    ]
  }
};

/**
 * Get example for a specific tool
 */
export function getToolExample(toolName) {
  return toolExamples[toolName] || null;
}

/**
 * Format example as JSON string for display
 */
export function formatExample(example) {
  if (!example || Object.keys(example).length === 0) {
    return '// No parameters needed';
  }
  return JSON.stringify(example, null, 2);
}

/**
 * Get parameter type from schema
 */
export function getParameterType(schema) {
  if (schema.type === 'array') {
    if (schema.items?.type) {
      return `array of ${schema.items.type}s`;
    }
    return 'array';
  }

  if (schema.type === 'object') {
    return 'object (see example)';
  }

  if (schema.enum) {
    return `one of: ${schema.enum.join(', ')}`;
  }

  return schema.type || 'unknown';
}

/**
 * Validate parameter value against schema
 */
export function validateParameter(value, schema) {
  if (schema.type === 'string') {
    return typeof value === 'string';
  }

  if (schema.type === 'number') {
    return typeof value === 'number' && !isNaN(value);
  }

  if (schema.type === 'boolean') {
    return typeof value === 'boolean';
  }

  if (schema.type === 'array') {
    return Array.isArray(value);
  }

  if (schema.type === 'object') {
    return typeof value === 'object' && value !== null;
  }

  return true; // Unknown type, allow it
}