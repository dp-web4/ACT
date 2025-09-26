#!/usr/bin/env node
import 'dotenv/config';  // Load environment variables
import express from 'express';
import cors from 'cors';
import { fileURLToPath } from 'url';
import path from 'path';
import MCPDiscovery from './discover.js';
import MCPConnector from './connector.js';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.static(path.join(__dirname, '..', 'public')));

// Serve tool-helpers.js from src directory
app.get('/tool-helpers.js', (req, res) => {
  res.sendFile(path.join(__dirname, 'tool-helpers.js'));
});

// Serve intelligent-helper.js from src directory
app.get('/intelligent-helper.js', (req, res) => {
  res.sendFile(path.join(__dirname, 'intelligent-helper.js'));
});

// Serve server-info.js from src directory
app.get('/server-info.js', (req, res) => {
  res.sendFile(path.join(__dirname, 'server-info.js'));
});

// Store active connections
const connections = new Map();

// API Routes

// Discover available MCP servers
app.get('/api/discover', async (req, res) => {
  try {
    const discovery = new MCPDiscovery();
    const servers = await discovery.discover();

    // Add connection status to each server
    const serversWithStatus = servers.map(server => ({
      ...server,
      connected: connections.has(server.name),
      tools: connections.get(server.name)?.tools || []
    }));

    res.json(serversWithStatus);
  } catch (error) {
    console.error('Discovery error:', error);
    res.status(500).json({ error: error.message });
  }
});

// Connect to a specific MCP server
app.post('/api/connect', async (req, res) => {
  try {
    const server = req.body;

    // Close existing connection if any
    if (connections.has(server.name)) {
      const existing = connections.get(server.name);
      await existing.disconnect();
      connections.delete(server.name);
    }

    // Create new connection
    const connector = new MCPConnector(server);
    await connector.connect();

    // Get available tools
    const tools = await connector.listTools();

    // Store connection
    connections.set(server.name, connector);

    res.json(tools);
  } catch (error) {
    console.error('Connection error:', error);
    res.status(500).json({ error: error.message });
  }
});

// Invoke a tool on a connected server
app.post('/api/invoke', async (req, res) => {
  try {
    const { server, tool, parameters } = req.body;
    console.log(`Invoking ${tool} on ${server} with:`, parameters);

    const connector = connections.get(server);
    if (!connector) {
      throw new Error(`Not connected to server: ${server}`);
    }

    const result = await connector.invokeTool(tool, parameters);
    console.log(`Result from ${tool}:`, result);
    res.json(result);
  } catch (error) {
    console.error('Invocation error:', error);
    res.status(500).json({ error: error.message });
  }
});

// Disconnect from a server
app.post('/api/disconnect', async (req, res) => {
  try {
    const { server } = req.body;

    if (connections.has(server)) {
      const connector = connections.get(server);
      await connector.disconnect();
      connections.delete(server);
    }

    res.json({ success: true });
  } catch (error) {
    console.error('Disconnect error:', error);
    res.status(500).json({ error: error.message });
  }
});

// Get server status
app.get('/api/status/:server', (req, res) => {
  const { server } = req.params;
  const connector = connections.get(server);

  if (!connector) {
    res.json({ connected: false });
  } else {
    res.json({
      connected: true,
      tools: connector.tools || [],
      metadata: connector.metadata || {}
    });
  }
});

// Start server
app.listen(PORT, () => {
  console.log(`
╔════════════════════════════════════════╗
║     ACT - Agentic Context Tool         ║
║     Direct MCP Server Interface        ║
╠════════════════════════════════════════╣
║  Server running on:                    ║
║  http://localhost:${PORT}                  ║
╠════════════════════════════════════════╣
║  Features:                             ║
║  • Discover MCP servers                ║
║  • Connect without AI intermediary     ║
║  • Direct tool invocation              ║
║  • Real-time results                   ║
╚════════════════════════════════════════╝
  `);
});