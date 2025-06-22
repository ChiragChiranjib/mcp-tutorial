# MCP Server Configuration Setup

This directory contains pre-configured JSON files for easy setup of all MCP transport modes.

## Quick Setup

### For Claude Desktop

1. **Copy the configuration:**
   ```bash
   cp claude_mcp.json ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

2. **Update the stdio path** in the copied file to match your actual installation path
3. **Restart Claude Desktop**

### For Cursor IDE

1. **Copy the configuration:**
   ```bash
   cp cursor_mcp.json ~/.cursor/mcp_servers.json
   ```

2. **Update the stdio path** in the copied file to match your actual installation path
3. **Restart Cursor**

## Before Using

### 1. Build All Binaries
```bash
mkdir -p bin
go build -o bin/stdio cmd/stdio/main.go
go build -o bin/sse cmd/sse/main.go
go build -o bin/streamable_http cmd/streamable_http/main.go
chmod +x bin/*
```

### 2. Start HTTP-based Servers
```bash
# Terminal 1: Start SSE server (port 8080)
./bin/sse

# Terminal 2: Start HTTP server (port 8081)
./bin/streamable_http
```

### 3. Update Paths
Edit the configuration files and replace:
```
/Users/chirag.chiranjib/razorpay/golang/mcp/tutorial/bin/stdio
```
With your actual absolute path to the stdio binary.

## Available Servers

All three servers provide identical functionality:

- **Tools:** `calculator`, `system_info`
- **Prompts:** `math_tutor`, `code_review`  
- **Resources:** `system://status`, `math://constants`

### Transport Methods

1. **tutorial-mcp-stdio** - Standard input/output (always available)
2. **tutorial-mcp-sse** - Server-Sent Events over HTTP (requires server running)
3. **tutorial-mcp-http** - HTTP with streaming (requires server running)

## Testing

After setup, test with any of these commands in Claude/Cursor:
- *"Calculate 25 * 4 using the calculator tool"*
- *"Show me system information"*
- *"Use the math tutor prompt for algebra"*

## Troubleshooting

- **Stdio not working:** Check that the binary path is absolute and executable
- **HTTP servers not working:** Ensure the servers are running and ports are available
- **Changes not taking effect:** Restart the application completely
- **Port conflicts:** Modify the port numbers in the server startup and configuration files 