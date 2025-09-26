import { spawn } from 'child_process';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { StdioClientTransport } from '@modelcontextprotocol/sdk/client/stdio.js';

/**
 * MCP Server Connector
 * Handles connection and communication with MCP servers
 */
class MCPConnector {
  constructor(serverConfig) {
    this.config = serverConfig;
    this.client = null;
    this.transport = null;
    this.tools = [];
    this.metadata = {};
    this.connected = false;
  }

  async connect() {
    try {
      // Prepare command and args
      let command, args;

      if (this.config.type === 'npx') {
        command = 'npx';
        args = [this.config.package, ...this.config.args];
      } else if (this.config.command) {
        const parts = this.config.command.split(' ');
        command = parts[0];
        args = parts.slice(1);
      } else {
        throw new Error('No command specified for server');
      }

      // Create the transport
      this.transport = new StdioClientTransport({
        command,
        args,
        env: { ...process.env, ...this.config.env }
      });

      // Create the client
      this.client = new Client({
        name: 'act-tool',
        version: '0.1.0'
      }, {
        capabilities: {}
      });

      // Connect
      await this.client.connect(this.transport);
      this.connected = true;

      // Get server info
      const serverInfo = await this.client.getServerInfo();
      this.metadata = serverInfo;

      // Get available tools
      await this.refreshTools();

      return true;
    } catch (error) {
      console.error('Connection error:', error);
      this.connected = false;
      throw error;
    }
  }

  async refreshTools() {
    if (!this.connected || !this.client) {
      throw new Error('Not connected to server');
    }

    try {
      const toolsResponse = await this.client.listTools();
      this.tools = toolsResponse.tools || [];
      return this.tools;
    } catch (error) {
      console.error('Error fetching tools:', error);
      // Some servers might not support tools - that's okay
      this.tools = [];
      return [];
    }
  }

  async listTools() {
    if (!this.connected) {
      await this.connect();
    }
    return this.tools;
  }

  async invokeTool(toolName, parameters = {}) {
    if (!this.connected || !this.client) {
      throw new Error('Not connected to server');
    }

    const tool = this.tools.find(t => t.name === toolName);
    if (!tool) {
      throw new Error(`Tool not found: ${toolName}`);
    }

    try {
      const result = await this.client.callTool({
        name: toolName,
        arguments: parameters
      });

      return {
        success: true,
        result: result.content,
        toolName,
        parameters
      };
    } catch (error) {
      console.error('Tool invocation error:', error);
      throw error;
    }
  }

  async disconnect() {
    if (this.client) {
      try {
        await this.client.close();
      } catch (error) {
        console.error('Error closing client:', error);
      }
    }

    if (this.transport) {
      try {
        await this.transport.close();
      } catch (error) {
        console.error('Error closing transport:', error);
      }
    }

    this.connected = false;
    this.client = null;
    this.transport = null;
    this.tools = [];
  }

  isConnected() {
    return this.connected;
  }
}

// Alternative simpler connector for testing without full MCP SDK
class SimpleMCPConnector {
  constructor(serverConfig) {
    this.config = serverConfig;
    this.process = null;
    this.tools = [];
    this.connected = false;
    this.messageId = 1;
    this.pendingRequests = new Map();
  }

  async connect() {
    return new Promise((resolve, reject) => {
      // Prepare command
      let command, args;

      if (this.config.type === 'npx') {
        command = 'npx';
        // Special handling for filesystem server
        if (this.config.name === 'filesystem') {
          // Filesystem server expects the path as the last argument after the package name
          args = [this.config.package, ...this.config.args];
        } else {
          args = [this.config.package, ...this.config.args];
        }
      } else if (this.config.command) {
        const parts = this.config.command.split(' ');
        command = parts[0];
        args = parts.slice(1);
      } else {
        reject(new Error('No command specified'));
        return;
      }

      // Spawn the process
      this.process = spawn(command, args, {
        env: { ...process.env, ...this.config.env },
        shell: true
      });

      let initTimeout = setTimeout(() => {
        reject(new Error('Connection timeout'));
        this.disconnect();
      }, 15000);  // Increased timeout for slower servers

      // Buffer for incomplete messages
      let buffer = '';

      // Handle stdout (JSON-RPC responses)
      this.process.stdout.on('data', (data) => {
        clearTimeout(initTimeout);

        buffer += data.toString();

        // Process complete messages (newline-delimited JSON)
        const lines = buffer.split('\n');
        buffer = lines.pop() || '';  // Keep incomplete line in buffer

        for (const line of lines) {
          if (line.trim()) {
            try {
              const message = JSON.parse(line);
              this.handleMessage(message);
            } catch (error) {
              // Some servers send non-JSON output, ignore it
              console.debug('Non-JSON output:', line);
            }
          }
        }
      });

      this.process.stderr.on('data', (data) => {
        // Many MCP servers output status to stderr, not all of it is errors
        const message = data.toString();
        if (message.toLowerCase().includes('error')) {
          console.error('Server error:', message);
        } else {
          console.debug('Server message:', message);
        }
      });

      this.process.on('error', (error) => {
        console.error('Process error:', error);
        reject(error);
      });

      this.process.on('close', (code) => {
        this.connected = false;
        console.log('Server process closed with code:', code);
      });

      // Wait a moment for the server to start up
      setTimeout(() => {
        console.log(`Initializing ${this.config.name} server...`);

        // Send initialization request (MCP standard format)
        this.sendRequest('initialize', {
          protocolVersion: '2024-11-05',  // Updated to current MCP version
          capabilities: {
            roots: {
              listChanged: true
            },
            sampling: {}
          },
          clientInfo: {
            name: 'act-tool',
            version: '0.1.0'
          }
        }).then((initResponse) => {
          console.log(`${this.config.name} initialized:`, initResponse);
          this.connected = true;
          clearTimeout(initTimeout);

          // Send initialized notification (required by MCP)
          this.sendNotification('notifications/initialized', {});

          // Try to list tools
          this.sendRequest('tools/list', {}).then(response => {
            console.log(`${this.config.name} tools:`, response.tools?.length || 0);
            this.tools = response.tools || [];
            resolve(true);
          }).catch((err) => {
            console.log(`${this.config.name} doesn't support tools:`, err.message);
            // Server might not support tools
            this.tools = [];
            resolve(true);
          });
        }).catch((err) => {
          console.error(`Failed to initialize ${this.config.name}:`, err);
          reject(err);
        });
      }, 500);  // Give server 500ms to start up
    });
  }

  handleMessage(message) {
    if (message.id && this.pendingRequests.has(message.id)) {
      const { resolve, reject } = this.pendingRequests.get(message.id);
      this.pendingRequests.delete(message.id);

      if (message.error) {
        reject(new Error(message.error.message));
      } else {
        resolve(message.result);
      }
    }
  }

  sendRequest(method, params = {}) {
    return new Promise((resolve, reject) => {
      const id = this.messageId++;
      const request = {
        jsonrpc: '2.0',
        id,
        method,
        params
      };

      this.pendingRequests.set(id, { resolve, reject });

      // Send the request
      this.process.stdin.write(JSON.stringify(request) + '\n');

      // Timeout after 30 seconds
      setTimeout(() => {
        if (this.pendingRequests.has(id)) {
          this.pendingRequests.delete(id);
          reject(new Error('Request timeout'));
        }
      }, 30000);
    });
  }

  sendNotification(method, params = {}) {
    const notification = {
      jsonrpc: '2.0',
      method,
      params
    };

    // Send the notification (no response expected)
    if (this.process && this.process.stdin) {
      this.process.stdin.write(JSON.stringify(notification) + '\n');
    }
  }

  async listTools() {
    return this.tools;
  }

  async invokeTool(toolName, parameters = {}) {
    if (!this.connected) {
      throw new Error('Not connected to server');
    }

    try {
      // Standard MCP tool call format
      const result = await this.sendRequest('tools/call', {
        name: toolName,
        arguments: parameters
      });

      return {
        success: true,
        result: result.content || result,  // Handle different response formats
        toolName,
        parameters
      };
    } catch (error) {
      // If tools/call fails, try the alternate format
      try {
        const result = await this.sendRequest(`${toolName}`, parameters);
        return {
          success: true,
          result: result.content || result,
          toolName,
          parameters
        };
      } catch (innerError) {
        throw error;  // Throw original error
      }
    }
  }

  async disconnect() {
    // Just close our connection, don't kill the process
    // Other applications might be using the same MCP server
    if (this.process && this.process.stdin) {
      try {
        // Close stdin to signal we're done
        this.process.stdin.end();
      } catch (error) {
        console.error('Error closing stdin:', error);
      }
    }

    // Clean up our references
    this.process = null;
    this.connected = false;
    this.tools = [];
    this.pendingRequests.clear();
  }

  isConnected() {
    return this.connected;
  }
}

// Export the simpler connector for now as MCP SDK might need more setup
export default SimpleMCPConnector;