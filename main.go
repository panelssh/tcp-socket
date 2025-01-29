// Package main provides a secure TCP server that executes shell commands
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

// Request represents the incoming command request structure
type Request struct {
	SecretKey string `json:"secretKey"` // Authentication key
	Command   string `json:"command"`   // Command to execute
}

// Response represents the server's response structure
type Response struct {
	Success bool   `json:"success"` // Indicates if the command executed successfully
	Message string `json:"message"` // Status or error message
	Output  string `json:"output"`  // Command execution output
}

// Configuration holds the server configuration
type Configuration struct {
	Host           string
	Port           string
	SecretKey      string
	AllowedAddress string
}

// getEnv retrieves an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// isIPAllowed checks if the given IP is in the allowed list
func isIPAllowed(ip string, allowedList string) bool {
	if allowedList == "%" {
		return true
	}
	allowedIPs := strings.Split(allowedList, ",")
	for _, allowed := range allowedIPs {
		if allowed == ip {
			return true
		}
	}
	return false
}

// sendResponse sends a JSON response back to the client
func sendResponse(conn net.Conn, response Response) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}

// loadConfiguration loads server configuration from flags and environment variables
func loadConfiguration() Configuration {
	config := Configuration{
		Host:           getEnv("HOST", "0.0.0.0"),
		Port:           getEnv("PORT", "3000"),
		SecretKey:      getEnv("SECRET_KEY", ""),
		AllowedAddress: getEnv("ALLOWED_ADDRESS", "%"),
	}

	flag.StringVar(&config.Host, "host", config.Host, "Bind Host/IP Address")
	flag.StringVar(&config.Port, "port", config.Port, "Listening Port")
	flag.StringVar(&config.SecretKey, "secret-key", config.SecretKey, "Secret Key")
	flag.StringVar(&config.AllowedAddress, "allowed-address", config.AllowedAddress, "List Allowed IP Address")
	flag.Parse()

	return config
}

// handleConnection processes incoming client connections
func handleConnection(conn net.Conn, config Configuration) {
	defer conn.Close()

	response := Response{Success: false, Message: "OK"}

	// Check IP allowlist
	clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
	if !isIPAllowed(clientIP, config.AllowedAddress) {
		response.Message = "Request not allowed"
		sendResponse(conn, response)
		return
	}

	// Read request
	buffer := make([]byte, 4096)
	read, err := conn.Read(buffer)
	if err != nil || read == 0 {
		response.Message = "Invalid request"
		sendResponse(conn, response)
		return
	}

	// Parse request
	var request Request
	if err := json.Unmarshal(buffer[:read], &request); err != nil {
		response.Message = "Invalid request format"
		sendResponse(conn, response)
		return
	}

	// Validate request
	if err := validateRequest(request, config.SecretKey); err != nil {
		response.Message = err.Error()
		sendResponse(conn, response)
		return
	}

	// Execute command
	output, err := executeCommand(request.Command)
	if err != nil {
		response.Message = err.Error()
		sendResponse(conn, response)
		return
	}

	response.Success = true
	response.Output = strings.TrimSpace(output)
	sendResponse(conn, response)
}

// validateRequest validates the incoming request
func validateRequest(request Request, secretKey string) error {
	if request.SecretKey == "" {
		return fmt.Errorf("secret key is required")
	}
	if request.Command == "" {
		return fmt.Errorf("command is required")
	}
	if request.SecretKey != secretKey {
		return fmt.Errorf("invalid secret key")
	}
	return nil
}

// executeCommand executes a shell command and returns its output
func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}
	return string(output), nil
}

func main() {
	config := loadConfiguration()

	if config.SecretKey == "" {
		fmt.Println("Error: Secret key is required")
		os.Exit(1)
	}

	// Start server
	listener, err := net.Listen("tcp", config.Host+":"+config.Port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[*] Server listening on %s\n", listener.Addr())

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn, config)
	}
}
