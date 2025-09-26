# ACT Tool - Direct MCP Server Interface

## 🚧 Experimental Status

This is an **experimental prototype** for exploring direct human interaction with Model Context Protocol (MCP) servers. The tool is in active discovery phase and serves as a learning platform for understanding MCP capabilities and requirements.

## Overview

The ACT (Agentic Context Tool) provides a web-based interface for:
- **Discovering** available MCP servers
- **Connecting** directly to MCP servers without AI intermediaries
- **Invoking** MCP tools through a user-friendly interface
- **Exploring** MCP server capabilities and data formats

## Features

### Current Capabilities
- 🔍 **Auto-discovery** of MCP servers (npx packages and local configurations)
- 🔌 **Direct connection** management without killing shared server processes
- 📝 **Intelligent input helpers** that analyze MCP schemas to provide examples
- ℹ️ **Server information** with descriptions, limitations, and usage tips
- 🎯 **Tool filtering** to hide deprecated tools
- 📋 **Copy-to-clipboard** for examples and results

### Supported MCP Servers
- **Filesystem** - File and directory operations (restricted to `/tmp` for safety)
- **Git** - Repository operations
- **Weather** - US weather forecasts and alerts
- **Memory** - Knowledge graph storage (ephemeral with npx)
- **Puppeteer** - Browser automation
- **PostgreSQL** - Database operations

## Installation & Usage

### Prerequisites
- Node.js 20+
- npm or yarn
- Environment variables (optional):
  - `WEATHER_API_KEY` for weather server
  - `DATABASE_URL` for PostgreSQL server

### Quick Start

```bash
# Install dependencies
npm install

# Start the server
npm start

# Open in browser
# Navigate to http://localhost:3000
```

### Usage Tips

1. **Discovery Phase**: Click "Discover Servers" to find available MCP servers
2. **Server Info**: Click the ? button next to each server for details
3. **Connect**: Select a server to connect and view its available tools
4. **Tool Usage**:
   - Review the example inputs provided
   - Fill in the parameter fields
   - Click "Invoke" to execute the tool
   - Results appear in the right panel

### Important Notes

#### Path Restrictions
The filesystem server is restricted to `/tmp` directory for security. Always use paths like:
- ✅ `/tmp/myfile.txt`
- ✅ `/tmp/myproject/`
- ❌ `~/Documents/file.txt` (not allowed)

#### Ephemeral Storage
When using `npx` to run MCP servers, data is **not persistent** between sessions. Each new connection creates a fresh instance without previous state.

#### Process Management
The tool uses intelligent process management:
- Discovery finds available servers without spawning them
- Connect creates connections (spawns if needed)
- Disconnect closes connections without killing the server process
- Multiple applications can share the same MCP server

## Architecture

### Components
- **Express Backend** (`src/server.js`) - REST API for MCP operations
- **MCP Connector** (`src/connector.js`) - Handles JSON-RPC communication
- **Discovery System** (`src/discover.js`) - Finds available MCP servers
- **Web Frontend** (`public/index.html`) - Three-panel UI
- **Intelligent Helpers** (`src/intelligent-helper.js`) - Schema-based example generation
- **Server Info** (`src/server-info.js`) - Detailed server documentation

### API Endpoints
- `GET /api/discover` - Find available MCP servers
- `POST /api/connect` - Connect to a specific server
- `POST /api/invoke` - Execute a tool on connected server
- `POST /api/disconnect` - Disconnect from server
- `GET /api/status/:server` - Get server connection status

## Development Status

### ✅ Completed
- Basic MCP server discovery and connection
- Tool invocation with parameter handling
- Intelligent input helpers
- Server information system
- Process management fixes

### 🚧 In Progress
- Enhanced error handling
- Better UI/UX for complex tool parameters
- Support for more MCP servers
- Persistent storage options

### 📋 Planned
- Tool result visualization
- Batch operations
- Configuration management
- Advanced filtering and search

## Known Limitations

1. **Browser Security**: Some operations may be limited by browser security policies
2. **Complex Parameters**: Tools requiring complex nested JSON may need manual formatting
3. **Timeout Issues**: Some servers may timeout during discovery
4. **Platform Specific**: Tested primarily on WSL/Linux environments

## Contributing

This is an experimental tool in active development. Key areas for contribution:
- Testing with additional MCP servers
- UI/UX improvements
- Documentation of MCP tool patterns
- Error handling enhancements

## License

See main ACT repository for license information.

## Acknowledgments

Built to explore the Model Context Protocol (MCP) ecosystem and enable direct human interaction with MCP servers for learning and experimentation.