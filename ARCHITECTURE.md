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

## Comparison of Implementations

| Feature | SSE | STDIO | Streamable HTTP |
|---------|-----|-------|----------------|
| **Transport** | HTTP + SSE | stdin/stdout | HTTP Streaming |
| **Port** | 8080 | N/A | 8081 |
| **State** | Stateful | Stateful | Stateless |
| **Keep-Alive** | Yes (10s) | N/A | No |
| **Use Case** | Web browsers, real-time apps | CLI tools, process integration | REST APIs, microservices |
| **Scalability** | Moderate | Single process | High (stateless) |
| **Complexity** | Medium | Low | Medium |

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