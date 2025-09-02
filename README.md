# Net-Cat Clone 🚀

A full-featured TCP chat server implementation in Go, recreating the functionality of NetCat in a Server-Client Architecture. This project demonstrates concurrent programming, network protocols, and real-time communication systems.

## 📋 Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Technical Implementation](#technical-implementation)
- [Learning Objectives](#learning-objectives)
- [Testing](#testing)
- [Logging System](#logging-system)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Documentation](#documentation)

## 🎯 Overview

This project implements a TCP-based chat server that can handle multiple concurrent client connections. It recreates NetCat functionality with enhanced features including user authentication, message broadcasting, chat persistence, and comprehensive logging.

### Key Capabilities
- **Server Mode**: Listen on specified ports for incoming connections
- **Multi-Client Support**: Handle up to 10 concurrent connections
- **Real-time Messaging**: Instant message broadcasting to all connected clients
- **User Management**: Username validation and collision detection
- **Chat Persistence**: Message history saved and restored
- **Comprehensive Logging**: Detailed server activity tracking

## ✨ Features

### Core Functionality
- ✅ TCP server listening on configurable ports
- ✅ Concurrent client connection handling
- ✅ Real-time message broadcasting
- ✅ Username authentication and validation
- ✅ Connection limit management (max 10 clients)
- ✅ Graceful client disconnect handling

### Enhanced Features
- 🔧 Welcome message display (ASCII art support)
- 📝 Chat history persistence and retrieval
- 📊 Comprehensive logging system (Info, Error, Access logs)
- 🛡️ Input validation and sanitization
- ⚡ Thread-safe operations with mutex protection
- 🌐 Cross-platform network compatibility

## 🏗️ Architecture

```
┌─────────────────┐    TCP Connection    ┌──────────────────┐
│   Client 1      │◄──────────────────► │                  │
│   (nc/telnet)   │                     │                  │
└─────────────────┘                     │                  │
                                        │   Chat Server    │
┌─────────────────┐    TCP Connection    │                  │
│   Client 2      │◄──────────────────► │  - User Mgmt     │
│   (nc/telnet)   │                     │  - Broadcasting  │
└─────────────────┘                     │  - Logging       │
                                        │  - Persistence   │
┌─────────────────┐    TCP Connection    │                  │
│   Client N      │◄──────────────────► │                  │
│   (nc/telnet)   │                     │                  │
└─────────────────┘                     └──────────────────┘
```

## 🚀 Installation

### Prerequisites
- Go 1.19 or higher
- Network connectivity for testing
- Terminal access (Linux/macOS/Windows)

### Setup
```bash
# Clone the repository
git clone <your-repo-url>
cd netcat-clone

# Build the project
go build -o TCPChat .

# Or run directly
go run . [port]
```

## 💡 Usage

### Starting the Server

#### Default Port (8989)
```bash
go run .
# or
go run . 8989
```

#### Custom Port
```bash
go run . 3000
```

#### Using Built Binary
```bash
./TCPChat 8989
```

### Connecting Clients

#### Using NetCat
```bash
nc localhost 8989
# or for remote connections
nc <server-ip> 8989
```

#### Using Telnet
```bash
telnet localhost 8989
```

#### Using Other Tools
```bash
# Any TCP client can connect
socat - TCP:localhost:8989
```

### Sample Session

**Server Terminal:**
```bash
$ go run . 8989
Chat server listening on port 8989
Server is ready! Waiting for connections...
Local connection: nc localhost 8989
Remote connection: nc <server-ip> 8989
New client connected from 127.0.0.1:54321
User 'alice' joined the chat
New client connected from 127.0.0.1:54322
User 'bob' joined the chat
```

**Client 1 (Alice):**
```
Welcome to TCP Chat Server!
[ENTER YOUR NAME]: alice
bob has joined the chat...
[2024-01-15 14:30:25][alice]: Hello everyone!
[2024-01-15 14:30:27][bob]: Hi Alice! How are you?
[2024-01-15 14:30:30][alice]: 
```

**Client 2 (Bob):**
```
Welcome to TCP Chat Server!
[ENTER YOUR NAME]: bob
=== Chat History ===
[2024-01-15 14:30:25][alice]: Hello everyone!
=== End History ===

[2024-01-15 14:30:27][bob]: Hi Alice! How are you?

[2024-01-15 14:30:25][alice]: Hello everyone!
[2024-01-15 14:30:30][bob]: 
```

## 📁 Project Structure

```
netcat-clone/
├── main.go                 # Application entry point
├── server/
│   └── server.go          # Core server implementation
├── logs/                  # Auto-generated log directory
│   ├── info.log          # General server operations
│   ├── error.log         # Error tracking
│   └── access.log        # Connection events
├── savedChat.txt         # Persistent chat history
├── pingoin.txt          # Welcome message (ASCII art)
├── README.md            # This file
└── go.mod              # Go module definition
```

## 🔧 Technical Implementation

### Concurrency Model
- **Goroutines**: Each client connection runs in its own goroutine
- **Mutex Protection**: `sync.RWMutex` ensures thread-safe access to shared data
- **Channel-free Design**: Uses mutex-protected maps for client management

### Data Structures
```go
type ChatServer struct {
    users     map[string]bool      // Username -> Online status
    clients   map[net.Conn]string  // Connection -> Username mapping
    mutex     sync.RWMutex         // Thread-safe access control
    port      string               // Server port
    // Logging components
    infoLog   *log.Logger
    errorLog  *log.Logger  
    accessLog *log.Logger
}
```

### Network Protocol
- **Transport**: TCP (Transmission Control Protocol)
- **Port Range**: Configurable (default: 8989)
- **Message Format**: `[timestamp][username]: message`
- **Connection Limit**: 10 concurrent clients maximum

### Key Algorithms
1. **Client Authentication**: Username validation with collision detection
2. **Message Broadcasting**: Efficient one-to-many message distribution
3. **Connection Management**: Graceful handling of client joins/leaves
4. **State Synchronization**: Mutex-protected shared state updates

## 🎓 Learning Objectives

This project demonstrates proficiency in:

### Go Programming Concepts
- **Structures**: Complex data type manipulation and methods
- **Concurrency**: Goroutines for parallel processing
- **Synchronization**: Mutexes for thread-safe operations
- **Error Handling**: Comprehensive error management
- **File I/O**: Reading/writing files for persistence

### Networking Fundamentals
- **TCP/UDP Protocols**: Understanding connection-oriented communication
- **Socket Programming**: Low-level network connection handling
- **IP Addressing**: Working with IPv4/IPv6 addresses
- **Port Management**: Binding and listening on network ports

### System Design
- **Client-Server Architecture**: Designing scalable network applications
- **State Management**: Handling shared state in concurrent systems
- **Resource Management**: Proper cleanup and resource allocation
- **Logging Systems**: Comprehensive activity tracking

### Software Engineering
- **Code Organization**: Clean, maintainable code structure
- **Documentation**: Comprehensive code and API documentation
- **Testing**: Network application testing strategies
- **Deployment**: Cross-platform application deployment

## 🧪 Testing

### Local Testing
```bash
# Terminal 1: Start server
go run . 8989

# Terminal 2: Connect client 1
nc localhost 8989

# Terminal 3: Connect client 2  
nc localhost 8989
```

### Network Testing
```bash
# Server machine
go run . 8989

# Client machine (replace with actual server IP)
nc 192.168.1.100 8989
```

### Automated Testing
```bash
# Test connection limit
for i in {1..12}; do
    echo "Connecting client $i"
    nc localhost 8989 &
done
```

### Load Testing
```bash
# Stress test with multiple rapid connections
for i in {1..50}; do
    (echo "user$i"; sleep 1; echo "Hello from user$i"; sleep 2) | nc localhost 8989 &
done
```

## 📊 Logging System

### Log Files Generated
- **`./logs/info.log`**: Server operations, user management, statistics
- **`./logs/error.log`**: Error conditions, connection failures, I/O issues
- **`./logs/access.log`**: Authentication events, connection tracking

### Sample Log Entries
```bash
# Info Log
INFO	2024/01/15 14:30:15	Chat server started successfully on port 8989
INFO	2024/01/15 14:30:20	User 'alice' joined from 192.168.1.101:54321

# Error Log  
ERROR	2024/01/15 14:30:30	server.go:85	Failed to send message to bob: connection reset

# Access Log
ACCESS	2024/01/15 14:30:15	New connection from 192.168.1.101:54321
ACCESS	2024/01/15 14:30:16	Username 'alice' accepted for client 192.168.1.101:54321
```

## ⚙️ Configuration

### Environment Variables
```bash
export NETCAT_PORT=8989          # Default port
export NETCAT_MAX_CLIENTS=10     # Connection limit
export NETCAT_LOG_LEVEL=INFO     # Logging verbosity
```

### Configuration Files
- **`pingoin.txt`**: Custom welcome message/ASCII art
- **Connection limits**: Modify `MaxConnections` constant in `server.go`

### Firewall Configuration
```bash
# Linux (ufw)
sudo ufw allow 8989/tcp

# Linux (iptables)  
sudo iptables -A INPUT -p tcp --dport 8989 -j ACCEPT

# Check if port is accessible
netstat -tlnp | grep 8989
```

## 🔍 Troubleshooting

### Common Issues

#### "Port already in use"
```bash
# Find process using the port
lsof -i :8989
# Kill the process
kill -9 <PID>
```

#### "Connection refused"
- Check if server is running: `netstat -tlnp | grep 8989`
- Verify firewall settings
- Confirm correct IP address and port

#### "No welcome message displayed"  
- Ensure `pingoin.txt` exists in project directory
- Check file permissions: `ls -la pingoin.txt`
- Review error logs: `tail -f logs/error.log`

#### "Messages not broadcasting"
- Check client connections in access logs
- Verify no network errors in error logs
- Confirm multiple clients are properly connected

### Debug Mode
```bash
# Enable verbose logging
go run . 8989 --verbose

# Monitor logs in real-time
tail -f logs/info.log logs/error.log logs/access.log
```

## 📚 Documentation

### Official References
- [NetCat Wikipedia](https://fr.wikipedia.org/wiki/NC)
- [Go Network Programming](https://golang.org/pkg/net/)
- [TCP Protocol Specification](https://tools.ietf.org/html/rfc793)

### Additional Resources
- [Go Concurrency Patterns](https://golang.org/doc/effective_go.html#concurrency)
- [Network Programming with Go](https://jan.newmarch.name/golang/)
- [Building Distributed Applications](https://pragprog.com/titles/rggo/distributed-services-with-go/)

## 📄 License

This project is developed for educational purposes. Feel free to use, modify, and distribute according to your institution's guidelines.

## 🏆 Acknowledgments

Special thanks to the Go community and networking protocol documentation that made this implementation possible.

---

Author **Kalla Moussa Gueye**

**Made with ❤️ and Go** | *Building the future of network communication*

