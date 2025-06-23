# MCP Tutorial Server Architecture

This document provides visual diagrams to help understand the different MCP (Model Context Protocol) server implementations in this tutorial project.

## Overview

The tutorial provides three different MCP server implementations:
- **SSE (Server-Sent Events)**: HTTP-based with real-time streaming
- **STDIO**: Standard input/output based communication  
- **Streamable HTTP**: HTTP-based with stateless streaming

All implementations share the same core MCP functionality but use different transport mechanisms.

## SSE Implementation Architecture

The SSE implementation uses HTTP with Server-Sent Events for real-time bidirectional communication.

```mermaid
graph TB
    subgraph "SSE MCP Server (Port 8080)"
        A[main.go] --> B[MCP Server Core]
        B --> C[SSE Server Wrapper]
        C --> D[HTTP Server with SSE]
        
        subgraph "MCP Core Components"
            B --> E[Tools]
            B --> F[Prompts] 
            B --> G[Resources]
        end
        
        subgraph "Tools"
            E --> H[Calculator Tool]
            E --> I[System Info Tool]
        end
        
        subgraph "Prompts"
            F --> J[Math Tutor Prompt]
            F --> K[Code Review Prompt]
        end
        
        subgraph "Resources"
            G --> L[System Status Resource]
            G --> M[Math Constants Resource]
        end
    end
    
    subgraph "Client Side"
        N[MCP Client] --> O[HTTP/SSE Connection]
        O --> P[Keep-Alive Mechanism<br/>10s interval]
    end
    
    P <--> D
    
    style A fill:#e1f5fe
    style C fill:#f3e5f5
    style D fill:#e8f5e8
```

## STDIO Implementation Architecture

The STDIO implementation uses standard input/output streams for communication, ideal for process-based integration.

```mermaid
graph TB
    subgraph "STDIO MCP Server"
        A[main.go] --> B[MCP Server Core]
        B --> C[STDIO Server Wrapper]
        C --> D[stdin/stdout Interface]
        
        subgraph "MCP Core Components"
            B --> E[Tools]
            B --> F[Prompts]
            B --> G[Resources]
        end
        
        subgraph "Tools"
            E --> H[Calculator Tool]
            E --> I[System Info Tool]
        end
        
        subgraph "Prompts"
            F --> J[Math Tutor Prompt]
            F --> K[Code Review Prompt]
        end
        
        subgraph "Resources"
            G --> L[System Status Resource]
            G --> M[Math Constants Resource]
        end
        
        subgraph "I/O Streams"
            D --> N[os.Stdin]
            D --> O[os.Stdout]
            P[os.Stderr] --> Q[Logging Output]
        end
    end
    
    subgraph "Client Process"
        R[MCP Client Process] --> S[Process Communication]
        S --> T[stdin pipe]
        S --> U[stdout pipe]
    end
    
    T <--> N
    U <--> O
    Q -.-> V[External Logs]
    
    style A fill:#e1f5fe
    style C fill:#fff3e0
    style D fill:#e8f5e8
```

## Streamable HTTP Implementation Architecture

The Streamable HTTP implementation provides HTTP-based communication with stateless streaming capabilities.

```mermaid
graph TB
    subgraph "Streamable HTTP MCP Server (Port 8081)"
        A[main.go] --> B[MCP Server Core]
        B --> C[Streamable HTTP Server]
        C --> D[HTTP Server<br/>Stateless Mode]
        
        subgraph "MCP Core Components"
            B --> E[Tools]
            B --> F[Prompts]
            B --> G[Resources]
        end
        
        subgraph "Tools"
            E --> H[Calculator Tool]
            E --> I[System Info Tool]
        end
        
        subgraph "Prompts"  
            F --> J[Math Tutor Prompt]
            F --> K[Code Review Prompt]
        end
        
        subgraph "Resources"
            G --> L[System Status Resource]
            G --> M[Math Constants Resource]
        end
    end
    
    subgraph "Client Side"
        N[HTTP Client] --> O[HTTP Requests]
        O --> P[Streaming Responses]
        Q[Load Balancer] --> R[Multiple Instances]
    end
    
    P <--> D
    R -.-> D
    
    style A fill:#e1f5fe
    style C fill:#fce4ec
    style D fill:#e8f5e8
```

## Client-Server Communication Patterns

### STDIO Communication Flow

The STDIO implementation uses standard input/output streams for direct process communication.

```mermaid
sequenceDiagram
    participant C as MCP Client
    participant P as Process
    participant S as MCP Server (STDIO)
    
    Note over C,S: Process Initialization
    C->>P: Start Server Process
    P->>S: Launch main.go
    S->>P: Ready on stdin/stdout
    
    Note over C,S: MCP Handshake
    C->>S: {"jsonrpc": "2.0", "method": "initialize", ...}
    Note over S: Read from stdin
    S->>C: {"jsonrpc": "2.0", "result": {"capabilities": ...}}
    Note over C: Read from stdout
    
    Note over C,S: Tool Execution
    C->>S: {"method": "tools/call", "params": {"name": "calculator", ...}}
    Note over S: Process via stdin
    S->>S: Execute Calculator Tool
    S->>C: {"result": {"content": [{"type": "text", "text": "2 + 2 = 4"}]}}
    Note over C: Receive via stdout
    
    Note over C,S: Resource Access
    C->>S: {"method": "resources/read", "params": {"uri": "system://status"}}
    S->>S: Generate System Status
    S->>C: {"result": {"contents": [{"uri": "system://status", ...}]}}
    
    Note over C,S: Process Termination
    C->>P: Terminate Process
    P->>S: SIGTERM/SIGINT
    S->>P: Clean Shutdown
```

### SSE Communication Flow

The SSE implementation uses HTTP with Server-Sent Events for real-time bidirectional communication.

```mermaid
sequenceDiagram
    participant C as MCP Client
    participant H as HTTP Server (Port 8080)
    participant S as MCP Server (SSE)
    
    Note over C,S: SSE Connection Setup
    C->>H: GET /sse
    H->>S: Initialize SSE Connection
    S->>H: Generate session_id
    H->>C: 200 OK + SSE Headers<br/>data: {"session_id": "abc123"}
    
    Note over C,S: Session Established
    C->>C: Store session_id = "abc123"
    H->>C: SSE Keep-Alive (every 10s)<br/>data: {"type": "ping"}
    
    Note over C,S: MCP Operations via Message Endpoint
    C->>H: POST /message?session_id=abc123<br/>{"method": "initialize", ...}
    H->>C: 202 Accepted (Request Queued)
    H->>S: Route to Session Handler
    S->>S: Process Initialize Request
    S->>H: {"result": {"capabilities": ...}}
    H->>C: SSE Push<br/>data: {"result": {"capabilities": ...}}
    
    Note over C,S: Tool Calls
    C->>H: POST /message?session_id=abc123<br/>{"method": "tools/call", "params": {...}}
    H->>C: 202 Accepted (Request Queued)
    H->>S: Route to Session Handler
    S->>S: Execute Calculator Tool
    S->>H: {"result": {"content": [...]}}
    H->>C: SSE Push<br/>data: {"result": {"content": [...]}}
    
    Note over C,S: Server-Initiated Notifications (Optional)
    S->>H: Server Event/Notification
    H->>C: SSE Push<br/>data: {"type": "notification", ...}
    
    Note over C,S: Connection Cleanup
    C->>H: Close SSE Connection
    H->>S: Cleanup Session
    S->>S: Remove session_id
```

### Streamable HTTP Communication Flow

The Streamable HTTP implementation provides stateless HTTP communication with optional SSE upgrade for notifications.

```mermaid
sequenceDiagram
    participant C as MCP Client
    participant H as HTTP Server (Port 8081)
    participant S as MCP Server (HTTP)
    
    Note over C,S: Direct MCP Endpoint Access
    C->>H: POST /mcp<br/>{"method": "initialize", ...}
    H->>S: Process Request (Stateless)
    S->>S: Handle Initialize
    S->>H: {"result": {"capabilities": ...}}
    H->>C: 200 OK + JSON Response
    
    Note over C,S: Tool Execution (Stateless)
    C->>H: POST /mcp<br/>{"method": "tools/call", "params": {"name": "calculator", ...}}
    H->>S: Process Request (No Session)
    S->>S: Execute Calculator Tool
    S->>H: {"result": {"content": [{"type": "text", "text": "5 * 3 = 15"}]}}
    H->>C: 200 OK + Tool Result
    
    Note over C,S: Resource Access (Stateless)
    C->>H: POST /mcp<br/>{"method": "resources/read", "params": {"uri": "math://constants"}}
    H->>S: Process Request (No Session)
    S->>S: Generate Math Constants
    S->>H: {"result": {"contents": [...]}}
    H->>C: 200 OK + Resource Data
    
    Note over C,S: Optional SSE Upgrade for Notifications
    alt Server Needs to Send Notifications
        C->>H: GET /sse-upgrade
        H->>S: Initialize Notification Channel
        S->>H: Generate temp session for notifications
        H->>C: 200 OK + SSE Headers<br/>data: {"upgrade": "success"}
        
        Note over C,S: Notification Delivery
        S->>H: Server Notification
        H->>C: SSE Push<br/>data: {"type": "server_notification", ...}
        
        Note over C,S: Continue Regular /mcp Calls
        C->>H: POST /mcp<br/>{"method": "tools/list"}
        H->>S: Process Request (Still Stateless)
        S->>H: {"result": {"tools": [...]}}
        H->>C: 200 OK + Tools List
    end
```

## MCP Protocol Flow

This diagram shows the typical request-response flow for MCP operations across all implementations.

```mermaid
sequenceDiagram
    participant C as MCP Client
    participant S as MCP Server
    participant T as Tool Handler
    participant R as Resource Handler
    participant P as Prompt Handler

    Note over C,P: Initialization Phase
    C->>S: Initialize Connection
    S->>C: Server Capabilities
    
    Note over C,P: Tool Execution Flow
    C->>S: List Available Tools
    S->>C: [calculator, system_info]
    C->>S: Call Tool (calculator)
    S->>T: Execute Calculator
    T->>T: Perform Math Operation
    T->>S: Return Result
    S->>C: Tool Response
    
    Note over C,P: Resource Access Flow
    C->>S: List Resources
    S->>C: [system://status, math://constants]
    C->>S: Read Resource (system://status)
    S->>R: Get System Status
    R->>R: Generate Status JSON
    R->>S: Return Status
    S->>C: Resource Content
    
    Note over C,P: Prompt Usage Flow
    C->>S: Get Prompt (math_tutor)
    S->>P: Generate Math Tutor Prompt
    P->>P: Build Custom Prompt
    P->>S: Return Prompt Template
    S->>C: Prompt Content
```

## Getting Started

### SSE Server
```bash
go run cmd/sse/main.go
# Server starts on http://localhost:8080
```

### STDIO Server  
```bash
go run cmd/stdio/main.go
# Communicates via stdin/stdout
```

### Streamable HTTP Server
```bash
go run cmd/streamable_http/main.go  
# Server starts on http://localhost:8081
```

## Core MCP Components

All implementations share these components:

### Tools
- **Calculator**: Performs basic math operations (add, subtract, multiply, divide, power, sqrt)
- **System Info**: Provides current time/date in various formats

### Prompts
- **Math Tutor**: Comprehensive math tutoring with customizable topics and levels
- **Code Review**: Detailed code analysis with language-specific guidance

### Resources
- **System Status**: Server status and uptime information (JSON)
- **Math Constants**: Common mathematical constants (π, e, φ, √2) with descriptions

## Quick Start Examples

### STDIO Example
```bash
# Start the server process and communicate via pipes
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | ./bin/stdio
```

### SSE Example
```bash
# 1. Start server: ./bin/sse
# 2. Open SSE connection: curl -N http://localhost:8080/sse
# 3. Use session_id for operations: 
curl -X POST "http://localhost:8080/message?session_id=abc123" \
  -H "Content-Type: application/json" \
  -d '{"method":"tools/list","id":1}'
```

### Streamable HTTP Example
```bash
# 1. Start server: ./bin/streamable_http
# 2. Direct MCP calls (no session needed):
curl -X POST "http://localhost:8081/mcp" \
  -H "Content-Type: application/json" \
  -d '{"method":"tools/list","id":1}'
``` 