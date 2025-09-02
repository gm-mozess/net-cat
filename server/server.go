package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const MaxConnections = 10

// ChatServer represents a TCP chat server that manages multiple concurrent client connections.
// It maintains user sessions, handles message broadcasting, and provides chat persistence.
type ChatServer struct {
	users     map[string]bool      // Maps usernames to their online status (true = online)
	clients   map[net.Conn]string  // Maps network connections to their associated usernames
	mutex     sync.RWMutex         // Read-write mutex for thread-safe access to shared data
	port      string               // TCP port number the server listens on
	infoLog   *log.Logger          // Logger for informational messages
	errorLog  *log.Logger          // Logger for error messages
	accessLog *log.Logger          // Logger for access/connection logs
}

// NewChatServer creates and initializes a new ChatServer instance with logging.
// It takes a port string and returns a pointer to the configured server with loggers.
func NewChatServer(port string) *ChatServer {
	// Create log files with proper error handling
	infoFile := createLogFile("./logs/info.log")
	errorFile := createLogFile("./logs/error.log")
	accessFile := createLogFile("./logs/access.log")

	return &ChatServer{
		users:     make(map[string]bool),
		clients:   make(map[net.Conn]string),
		port:      port,
		infoLog:   log.New(infoFile, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog:  log.New(errorFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		accessLog: log.New(accessFile, "ACCESS\t", log.Ldate|log.Ltime),
	}
}

// createLogFile creates a log file and ensures the directory exists.
// If file creation fails, it falls back to stdout.
func createLogFile(filename string) *os.File {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("./logs", 0755); err != nil {
		fmt.Printf("Warning: Could not create logs directory: %v\n", err)
		return os.Stdout
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Warning: Could not create log file %s: %v\n", filename, err)
		return os.Stdout
	}
	return file
}

// Start begins listening for TCP connections on the configured port.
// It runs indefinitely, accepting new connections and spawning goroutines to handle them.
// Returns an error if the server fails to start listening.
func (cs *ChatServer) Start() error {
	// Create a TCP listener on the specified port
	listener, err := net.Listen("tcp", ":"+cs.port)
	if err != nil {
		cs.errorLog.Printf("Failed to start server on port %s: %v", cs.port, err)
		return fmt.Errorf("failed to start server: %w", err)
	}

	cs.infoLog.Printf("Chat server started successfully on port %s", cs.port)
	fmt.Printf("Chat server listening on port %s\n", cs.port)
	fmt.Println("Waiting for connections...")

	// Clear any existing chat history on server startup
	cs.emptyFile()
	
	// Ensure the listener is closed when the function exits
	defer listener.Close()

	// Accept connections in an infinite loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			cs.errorLog.Printf("Failed to accept connection: %v", err)
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		
		cs.accessLog.Printf("New connection from %s", conn.RemoteAddr())
		fmt.Printf("New client connected from %s\n", conn.RemoteAddr())
		
		// Handle each connection concurrently
		go cs.handleConnection(conn)
	}
}

// handleConnection manages the entire lifecycle of a client connection.
// It handles connection limits, welcome messages, user authentication, and message processing.
func (cs *ChatServer) handleConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	
	// Ensure connection cleanup and logging when function exits
	defer func() {
		conn.Close()
		cs.accessLog.Printf("Client %s disconnected", clientAddr)
		fmt.Printf("Client %s disconnected\n", clientAddr)
	}()

	// Check if server has reached maximum capacity
	cs.mutex.RLock()
	connectionCount := len(cs.clients)
	cs.mutex.RUnlock()

	if connectionCount >= MaxConnections {
		cs.infoLog.Printf("Connection rejected from %s: server full (%d/%d)", clientAddr, connectionCount, MaxConnections)
		fmt.Fprintln(conn, "Server is full. Please try again later.")
		return
	}

	// Send welcome message to the new client
	welcomeMsg, err := cs.getWelcomeMessage()
	if err != nil {
		cs.errorLog.Printf("Could not load welcome message: %v", err)
		fmt.Printf("Warning: Could not load welcome message: %v\n", err)
		fmt.Fprint(conn, "Welcome to the chat server!\n")
	} else {
		fmt.Fprint(conn, string(welcomeMsg))
		fmt.Fprint(conn, "\n")
	}

	// Obtain and validate username from client
	username := cs.getUserName(conn)
	if username == "" {
		cs.accessLog.Printf("Client %s failed to provide valid username", clientAddr)
		return // Connection was closed or username acquisition failed
	}

	cs.infoLog.Printf("User '%s' joined from %s", username, clientAddr)
	fmt.Printf("User '%s' joined the chat\n", username)

	// Register the new user in the server's data structures
	cs.mutex.Lock()
	cs.clients[conn] = username
	cs.users[username] = true
	currentUsers := len(cs.clients)
	cs.mutex.Unlock()

	cs.infoLog.Printf("Active users: %d/%d", currentUsers, MaxConnections)

	// Announce new user to existing clients and send chat history
	cs.broadcastJoin(conn, username)
	cs.sendChatHistory(conn)

	// Enter main message handling loop
	cs.handleMessages(conn, username)
}

// getUserName prompts the client for a username and validates it.
// It ensures the username is alphanumeric, within length limits, and not already taken.
// Returns the validated username in lowercase, or empty string if connection fails.
func (cs *ChatServer) getUserName(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	clientAddr := conn.RemoteAddr().String()

	for {
		fmt.Fprint(conn, "[ENTER YOUR NAME]: ")

		username, err := reader.ReadString('\n')
		if err != nil {
			cs.errorLog.Printf("Error reading username from %s: %v", clientAddr, err)
			return "" // Connection error or client disconnected
		}

		username = strings.TrimSpace(username)

		// Validate username format and length
		if !cs.isValidUsername(username) {
			cs.accessLog.Printf("Invalid username attempt from %s: '%s'", clientAddr, username)
			fmt.Fprintln(conn, "Username must be alphanumeric and max 20 characters!")
			continue
		}

		// Check if username is already in use
		cs.mutex.RLock()
		taken := cs.users[strings.ToLower(username)]
		cs.mutex.RUnlock()

		if taken {
			cs.accessLog.Printf("Username collision from %s: '%s' already taken", clientAddr, username)
			fmt.Fprintln(conn, "This username is already taken!")
			continue
		}

		cs.infoLog.Printf("Username '%s' accepted for client %s", username, clientAddr)
		return strings.ToLower(username)
	}
}

// isValidUsername checks if a username meets the server's requirements.
// Username must be 1-20 characters long and contain only alphanumeric characters.
func (cs *ChatServer) isValidUsername(username string) bool {
	if len(username) == 0 || len(username) > 20 {
		return false
	}

	for _, char := range username {
		if !((char >= '0' && char <= '9') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z')) {
			return false
		}
	}
	return true
}

// handleMessages manages the main chat loop for a connected client.
// It reads messages from the client, broadcasts them to other users, and saves them to file.
// The loop continues until the client disconnects.
func (cs *ChatServer) handleMessages(conn net.Conn, username string) {
	reader := bufio.NewReader(conn)
	messageCount := 0

	for {
		// Display message prompt with current timestamp
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(conn, "[%s][%s]: ", timestamp, username)

		// Read message from client
		message, err := reader.ReadString('\n')
		if err != nil {
			cs.accessLog.Printf("Client %s (%s) disconnected after %d messages", username, conn.RemoteAddr(), messageCount)
			break // Client disconnected or connection error
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue // Skip empty messages
		}

		messageCount++
		cs.infoLog.Printf("Message from %s: %s", username, message)

		// Format and broadcast the message to all other clients
		fullMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, username, message)
		cs.broadcastMessage(conn, fullMessage)
		cs.saveMessage(fullMessage)
	}

	// Clean up when client disconnects
	cs.handleDisconnection(conn, username)
}

// broadcastMessage sends a message to all connected clients except the sender.
// It creates a snapshot of current clients to avoid holding locks during network I/O.
func (cs *ChatServer) broadcastMessage(sender net.Conn, message string) {
	// Create a snapshot of current clients to avoid lock contention
	cs.mutex.RLock()
	clients := make(map[net.Conn]string)
	for conn, username := range cs.clients {
		clients[conn] = username
	}
	cs.mutex.RUnlock()

	broadcastCount := 0
	// Send message to all clients except the sender
	for conn, username := range clients {
		if conn != sender {
			_, err := fmt.Fprintf(conn, "\n%s\n", message)
			if err != nil {
				cs.errorLog.Printf("Failed to send message to %s: %v", username, err)
				continue
			}
			broadcastCount++
			
			// Redisplay the message prompt for each client
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(conn, "[%s][%s]: ", timestamp, username)
		}
	}
	
	cs.infoLog.Printf("Message broadcast to %d clients", broadcastCount)
}

// broadcastJoin announces when a new user joins the chat.
// It notifies all existing clients about the new user and redisplays their prompts.
func (cs *ChatServer) broadcastJoin(newConn net.Conn, username string) {
	joinMessage := fmt.Sprintf("%s has joined the chat...", username)

	// Get all existing clients (excluding the new user)
	cs.mutex.RLock()
	clients := make(map[net.Conn]string)
	for conn, user := range cs.clients {
		if conn != newConn {
			clients[conn] = user
		}
	}
	cs.mutex.RUnlock()

	notifyCount := 0
	// Announce the new user to existing clients
	for conn, user := range clients {
		_, err := fmt.Fprintf(conn, "\n%s\n", joinMessage)
		if err != nil {
			cs.errorLog.Printf("Failed to notify %s of new user: %v", user, err)
			continue
		}
		notifyCount++
		
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(conn, "[%s][%s]: ", timestamp, user)
	}
	
	cs.infoLog.Printf("Join notification sent to %d existing users", notifyCount)
}

// handleDisconnection cleans up when a client leaves the chat.
// It removes the client from server data structures and notifies remaining users.
func (cs *ChatServer) handleDisconnection(conn net.Conn, username string) {
	// Remove client from server data structures
	cs.mutex.Lock()
	delete(cs.clients, conn)
	delete(cs.users, username)

	// Create snapshot of remaining clients for notification
	clients := make(map[net.Conn]string)
	for c, u := range cs.clients {
		clients[c] = u
	}
	remainingUsers := len(cs.clients)
	cs.mutex.Unlock()

	cs.infoLog.Printf("User '%s' left the chat. Remaining users: %d/%d", username, remainingUsers, MaxConnections)

	// Notify remaining users about the departure
	leaveMessage := fmt.Sprintf("%s has left the chat...", username)
	notifyCount := 0
	
	for conn, user := range clients {
		_, err := fmt.Fprintf(conn, "\n%s\n", leaveMessage)
		if err != nil {
			cs.errorLog.Printf("Failed to notify %s of user departure: %v", user, err)
			continue
		}
		notifyCount++
		
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(conn, "[%s][%s]: ", timestamp, user)
	}
	
	cs.infoLog.Printf("Leave notification sent to %d remaining users", notifyCount)
}

// sendChatHistory sends the stored chat history to a newly connected client.
// This allows new users to see previous messages in the conversation.
func (cs *ChatServer) sendChatHistory(conn net.Conn) {
	file, err := os.Open("./savedChat.txt")
	if err != nil {
		cs.infoLog.Printf("No chat history available: %v", err)
		return // Chat history file doesn't exist or can't be opened
	}
	defer file.Close()

	lineCount := 0
	// Read and send each line of chat history
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(conn, scanner.Text())
		lineCount++
	}
	
	if lineCount > 0 {
		cs.infoLog.Printf("Sent %d lines of chat history to new user", lineCount)
	}
}

// saveMessage appends a chat message to the persistent chat history file.
// Messages are stored in "./savedChat.txt" for retrieval by new clients.
func (cs *ChatServer) saveMessage(message string) {
	file, err := os.OpenFile("./savedChat.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		cs.errorLog.Printf("Error saving message to file: %v", err)
		fmt.Printf("Error saving message: %v\n", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, message)
	if err != nil {
		cs.errorLog.Printf("Error writing message to file: %v", err)
	}
}

// emptyFile clears the chat history file when the server starts.
// This ensures each server session begins with a clean chat history.
func (cs *ChatServer) emptyFile() {
	file, err := os.OpenFile("./savedChat.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		cs.errorLog.Printf("Error clearing chat history: %v", err)
		fmt.Printf("Error clearing chat history: %v\n", err)
		return
	}
	file.Close()
	cs.infoLog.Println("Chat history cleared for new session")
}

// getWelcomeMessage reads and returns the welcome message from a file.
// The welcome message is typically ASCII art or server information displayed to new clients.
func (cs *ChatServer) getWelcomeMessage() ([]byte, error) {
	data, err := os.ReadFile("./pingoin.txt")
	if err != nil {
		return nil, fmt.Errorf("could not read welcome message file: %w", err)
	}
	cs.infoLog.Printf("Welcome message loaded successfully (%d bytes)", len(data))
	return data, nil
}