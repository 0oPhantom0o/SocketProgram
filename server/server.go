package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Configs struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Protocol string `yaml:"protocol"`
}

var conf = Configs{}
var mu sync.Mutex
var clientID string
var InfoLog *log.Logger

func init() {
	config()
	logData()
}
func main() {

	fmt.Printf("Server is listening on port %s...\n", conf.Port)

	listen()

}

func listen() {
	listener, err := net.Listen(conf.Protocol, conf.Host+conf.Port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr())
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)

	}
}

func handleClient(conn net.Conn) {
	for {

		buffer := make([]byte, 2000000)

		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}
		message := string(buffer[:bytesRead])
		message = replacement(message)
		InfoLog.Println(message + "\n")

		fmt.Printf("Received: %s", message)

		mu.Lock()

		clientID = conn.RemoteAddr().String()

		_, writeErr := conn.Write([]byte(clientID + " " + message))
		if writeErr != nil {
			log.Fatal("Error writing to client:", clientID, writeErr)
		}

		fmt.Println(clientID + " " + message)
		mu.Unlock()
	}
}
func config() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	fmt.Println("Raw YAML content:", string(data))
	// Unmarshal the YAML data into a Config struct
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	if conf.Port == "" || conf.Host == "" {
		log.Println("One or more fields are not populated.")
	}

}

func logData() {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(logFile, "Info "+time.Now().Format(time.RFC3339Nano)+clientID+"\n", log.Ldate|log.Ltime|log.Lshortfile)
}
func replacement(message string) string {
	message = strings.Replace(message, "#", "From:", -1)
	message = strings.Replace(message, "^", " To:", -1)
	message = strings.Replace(message, "$", " Amount:", -1)
	message = strings.Replace(message, "/", " Bank:", -1)
	message = strings.Replace(message, "-", " Time:", -1)
	//message = strings.Replace(message, "|", "", 1)

	parts := strings.Split(message, "|")

	// Initialize a new string to hold the result
	var result strings.Builder

	// Track packet count and iterate over each part
	packetCount := 1
	for _, part := range parts {
		if part != "" {
			result.WriteString(fmt.Sprintf("packet:%d %s ", packetCount, part))
			packetCount++
		}
	}

	// Trim any trailing spaces and print the result
	finalString := strings.TrimSpace(result.String())
	return finalString
}
