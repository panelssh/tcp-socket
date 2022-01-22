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

type Request struct {
	RequestId string `json:"request_id"`
	SecretKey string `json:"secret_key"`
	Command   string `json:"command"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output"`
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func sendResponse(conn net.Conn, response Response) {
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	defaultHost := getEnv("HOST", "0.0.0.0")
	defaultPort := getEnv("PORT", "3000")
	defaultSecretKey := getEnv("SECRET_KEY", "test")
	defaultAllowedAddress := getEnv("ALLOWED_ADDRESS", "%")

	Host := flag.String("host", defaultHost, "Bind Host/IP Address")
	Port := flag.String("port", defaultPort, "Listening Port")
	SecretKey := flag.String("secret-key", defaultSecretKey, "Secret Key")
	AllowedAddress := flag.String("allowed-address", defaultAllowedAddress, "List Allowed IP Address")

	flag.Parse()

	if *SecretKey == "" {
		fmt.Println("You must provide secret key")
		return
	}

	address := *Host + ":" + *Port
	listen, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("[*] Listening on %s\n", listen.Addr())

	for {
		response := Response{Success: false, Message: "OK", Output: ""}

		conn, err := listen.Accept()
		if err != nil {
			response.Message = err.Error()
			sendResponse(conn, response)
			continue
		}

		remoteAddress := strings.Split(conn.RemoteAddr().String(), ":")
		if *AllowedAddress != "%" {
			AllowedAddresses := strings.Split(*AllowedAddress, ",")
			if !stringInSlice(remoteAddress[0], AllowedAddresses) {
				response.Message = "Request Not Allowed"
				sendResponse(conn, response)
				continue
			}
		}

		bytes := make([]byte, 4096)
		read, err := conn.Read(bytes)
		if err != nil {
			response.Message = err.Error()
			sendResponse(conn, response)
			continue
		}

		if read == 0 {
			response.Message = "Empty Request"
			sendResponse(conn, response)
			continue
		}

		request := Request{}
		err = json.Unmarshal(bytes[:read], &request)
		if err != nil {
			response.Message = err.Error()
			sendResponse(conn, response)
			continue
		}

		if request.RequestId == "" {
			response.Message = "Request ID Is Required"
			sendResponse(conn, response)
			continue
		}

		if request.SecretKey == "" {
			response.Message = "Secret Key Is Required"
			sendResponse(conn, response)
			continue
		}

		if request.Command == "" {
			response.Message = "Command Is Required"
			sendResponse(conn, response)
			continue
		}

		if request.SecretKey != *SecretKey {
			response.Message = "Invalid Secret Key"
			sendResponse(conn, response)
			continue
		}

		cmd, err := exec.Command("bash", "-c", request.Command).Output()
		if err != nil {
			response.Message = err.Error()
			sendResponse(conn, response)
			continue
		}

		response.Success = true
		response.Output = strings.TrimSpace(string(cmd))
		sendResponse(conn, response)
	}
}
