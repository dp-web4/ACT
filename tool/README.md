# ACT Tool - Direct MCP Server Interface

A simple web-based tool that lets you discover and interact with MCP (Model Context Protocol) servers directly, without requiring an AI model as intermediary.

## Features

- **Server Discovery**: Automatically discovers available MCP servers (npx-based and local)
- **Direct Connection**: Connect to MCP servers without AI intermediary
- **Tool Exploration**: Browse available tools and their parameters
- **Direct Invocation**: Execute tools with custom parameters
- **Real-time Results**: See responses directly from the servers

## Quick Start

```bash
# Install dependencies
cd tool
npm install

# Start the server
npm start

# Or run in development mode with auto-reload
npm run dev
```

Then open http://localhost:3000 in your browser.

## How It Works

1. **Discovery Phase**: The tool scans for available MCP servers:
   - Known npx-runnable servers (filesystem, git, weather, etc.)
   - Local MCP server configurations in standard locations
   - Tests connectivity to each discovered server

2. **Connection Phase**: When you select a server:
   - Spawns the MCP server process
   - Establishes JSON-RPC communication
   - Retrieves available tools and their schemas

3. **Interaction Phase**: For each tool:
   - Displays parameters with appropriate input fields
   - Sends tool invocation requests
   - Shows results in real-time

## Architecture

```
Browser UI
    ↓
Express Server (port 3000)
    ↓
MCP Connector (SimpleMCPConnector)
    ↓
MCP Server Process (spawned)
```

## Components

### Frontend (`public/index.html`)
- Single-page application
- No framework dependencies
- Real-time server status updates
- Dynamic tool parameter forms

### Backend (`src/`)
- `server.js`: Express API server
- `discover.js`: MCP server discovery logic
- `connector.js`: MCP server connection management

### API Endpoints

- `GET /api/discover`: Find available MCP servers
- `POST /api/connect`: Connect to a specific server
- `POST /api/invoke`: Execute a tool on connected server
- `POST /api/disconnect`: Close server connection
- `GET /api/status/:server`: Get server connection status

## Supported MCP Servers

Currently discovers:
- `@modelcontextprotocol/server-filesystem`: File operations
- `@cyanheads/git-mcp-server`: Git operations
- `@h1deya/mcp-server-weather`: Weather information
- `@modelcontextprotocol/server-memory`: Knowledge graph
- `@modelcontextprotocol/server-puppeteer`: Browser automation
- `@modelcontextprotocol/server-postgres`: Database operations

## Configuration

Environment variables are loaded from `.env` file in the tool directory.

Create a `.env` file with:
```bash
# Weather API key for @h1deya/mcp-server-weather
WEATHER_API_KEY=your_api_key_here

# PostgreSQL connection string
DATABASE_URL=postgresql://localhost/test
```

Get your own Weather API key at: https://www.weatherapi.com/

## Development

The connector uses a simplified JSON-RPC implementation for MCP communication. Future versions will integrate the full MCP SDK for better compatibility.

### Adding New Servers

Edit `src/discover.js` and add to the `knownServers` array:

```javascript
{
  name: 'my-server',
  package: '@org/mcp-server-name',
  description: 'What it does',
  args: ['--any', '--args'],
  env: { ANY_ENV: 'value' }
}
```

## Limitations

- Currently uses simplified MCP connector (full SDK integration coming)
- Some MCP servers may require additional configuration
- Tool results are shown as raw JSON

## Next Steps

This is a starting point that can be expanded with:
- Full MCP SDK integration
- Tool result formatting and visualization
- Server configuration persistence
- Multi-server simultaneous connections
- Tool chaining and workflows
- Authentication and security features
- LCT integration for Web4 compliance

## License

Part of the ACT project - AGPL-3.0