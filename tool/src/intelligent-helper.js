/**
 * Intelligent Tool Helper Generator
 * Automatically generates examples and tips from MCP schema metadata
 * Context-agnostic and self-configuring
 */

export class IntelligentHelper {
  /**
   * Generate example value based on schema type and description
   */
  static generateExampleValue(schema, name = '') {
    const { type, description = '', enum: enumValues, items, properties } = schema;

    // Check description for hints
    const descLower = description.toLowerCase();

    switch (type) {
      case 'string':
        // Check for enums first
        if (enumValues && enumValues.length > 0) {
          return enumValues[0];
        }

        // Infer from name and description
        if (name.includes('path') || descLower.includes('path')) {
          // Filesystem server is restricted to /tmp
          return descLower.includes('directory') ? '/tmp/example_dir' : '/tmp/example.txt';
        }
        if (name.includes('name') || descLower.includes('name')) {
          return 'Example_Name';
        }
        if (name.includes('type') || descLower.includes('type')) {
          return 'example_type';
        }
        if (descLower.includes('state') && descLower.includes('code')) {
          return 'CA';
        }
        if (descLower.includes('url')) {
          return 'https://example.com';
        }
        if (name.includes('content') || descLower.includes('content')) {
          return 'Example content here';
        }
        if (name.includes('query') || descLower.includes('query') || descLower.includes('search')) {
          return 'search query';
        }
        if (name.includes('pattern') || descLower.includes('pattern')) {
          return '*.txt';
        }

        // Extract example from description if present
        const exampleMatch = descLower.match(/\(e\.g\.?\s*([^)]+)\)/);
        if (exampleMatch) {
          return exampleMatch[1].split(',')[0].trim();
        }

        return `example_${name || 'string'}`;

      case 'number':
        // Infer from name and description
        if (name.includes('latitude') || descLower.includes('latitude')) {
          return 37.7749;
        }
        if (name.includes('longitude') || descLower.includes('longitude')) {
          return -122.4194;
        }
        if (name.includes('days') || descLower.includes('days')) {
          return 3;
        }
        if (name.includes('count') || descLower.includes('count')) {
          return 10;
        }
        if (name.includes('limit') || descLower.includes('limit')) {
          return 100;
        }
        if (name.includes('offset') || descLower.includes('offset')) {
          return 0;
        }

        return 42;

      case 'boolean':
        return true;

      case 'array':
        if (items) {
          // Generate 1-2 example items
          const exampleItem = this.generateExampleValue(items, name.replace(/s$/, ''));
          return [exampleItem];
        }
        return ['item1', 'item2'];

      case 'object':
        if (properties) {
          const example = {};
          for (const [propName, propSchema] of Object.entries(properties)) {
            example[propName] = this.generateExampleValue(propSchema, propName);
          }
          return example;
        }
        return { key: 'value' };

      default:
        return null;
    }
  }

  /**
   * Generate tips from schema information
   */
  static generateTips(schema, toolName = '') {
    const tips = [];
    const { type, description = '', properties, required = [], enum: enumValues } = schema;

    // Add type-specific tips
    if (type === 'object' && properties) {
      const requiredFields = required.length > 0 ? required.join(', ') : 'none';
      tips.push(`Required fields: ${requiredFields}`);

      // Add property-specific tips
      for (const [propName, propSchema] of Object.entries(properties)) {
        if (propSchema.enum) {
          tips.push(`${propName} must be one of: ${propSchema.enum.join(', ')}`);
        }
        if (propSchema.type === 'array') {
          tips.push(`${propName} accepts multiple values`);
        }
      }
    }

    if (type === 'array' && schema.items) {
      tips.push('Provide an array of items');
      if (schema.items.type === 'object') {
        tips.push('Each item should be an object with the structure shown in the example');
      }
    }

    // Extract constraints from description
    const descLower = description.toLowerCase();
    if (descLower.includes('must be')) {
      const mustMatch = description.match(/must be ([^.]+)/i);
      if (mustMatch) {
        tips.push(mustMatch[0]);
      }
    }

    // Add format hints
    if (enumValues) {
      tips.push(`Allowed values: ${enumValues.join(', ')}`);
    }

    // Tool-specific tips based on common patterns
    if (toolName.includes('create') || toolName.includes('add')) {
      tips.push('Creates new data - ensure unique identifiers');
    }
    if (toolName.includes('delete') || toolName.includes('remove')) {
      tips.push('This action cannot be undone');
    }
    if (toolName.includes('search') || toolName.includes('find')) {
      tips.push('Supports partial matches unless specified otherwise');
    }

    // Filesystem-specific tips
    if (toolName.includes('file') || toolName.includes('directory')) {
      tips.push('⚠️ Filesystem restricted to /tmp directory for security');
      tips.push('Use paths like: /tmp/myfile.txt or /tmp/mydir/');
    }

    return tips;
  }

  /**
   * Generate complete helper for a tool based on its schema
   */
  static generateToolHelper(tool) {
    const { name, description, inputSchema } = tool;

    if (!inputSchema) {
      return {
        description: description || 'No description available',
        example: {},
        tips: ['No parameters required']
      };
    }

    // Generate example from schema
    const example = this.generateExampleValue(inputSchema, name);

    // Generate tips
    const tips = this.generateTips(inputSchema, name);

    // Add description-based tips
    if (description) {
      const descTips = this.extractTipsFromDescription(description);
      tips.push(...descTips);
    }

    return {
      description: description || 'No description available',
      example: example || {},
      tips: tips.length > 0 ? tips : ['Check the schema for parameter details']
    };
  }

  /**
   * Extract additional tips from tool description
   */
  static extractTipsFromDescription(description) {
    const tips = [];
    const descLower = description.toLowerCase();

    // Location-specific tips
    if (descLower.includes('us') || descLower.includes('united states')) {
      tips.push('US locations only');
    }

    // Format tips
    if (descLower.includes('two-letter') || descLower.includes('2-letter')) {
      tips.push('Use two-letter codes');
    }

    // Action tips
    if (descLower.includes('recursive')) {
      tips.push('Operates recursively on nested structures');
    }

    return tips;
  }

  /**
   * Generate helpers for all tools from a server
   */
  static generateAllHelpers(tools) {
    const helpers = {};

    for (const tool of tools) {
      helpers[tool.name] = this.generateToolHelper(tool);
    }

    return helpers;
  }

  /**
   * Infer the input widget type from schema
   */
  static inferWidgetType(schema) {
    const { type, properties, enum: enumValues } = schema;

    // Simple cases
    if (enumValues) {
      return 'select';
    }

    if (type === 'boolean') {
      return 'checkbox';
    }

    if (type === 'number') {
      return 'number';
    }

    if (type === 'string') {
      const name = schema.name || '';
      const desc = (schema.description || '').toLowerCase();

      if (name.includes('content') || desc.includes('content') || desc.includes('text')) {
        return 'textarea';
      }
      return 'text';
    }

    if (type === 'array' || (type === 'object' && properties)) {
      return 'json';
    }

    return 'text';
  }
}

// Export for use in browser
if (typeof window !== 'undefined') {
  window.IntelligentHelper = IntelligentHelper;
}