package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	NetWorkIp = "0.0.0.0"
	port      = ":9009"
	Protocol  = "tcp"
)

type SmS struct {
	From          string
	To            string
	Amount        int
	Bank          string
	OperationTime time.Time
}

var mu sync.Mutex

func main() {

	fmt.Printf("Server is listening on port %s...\n", port)

	// count client and decline new user if we had 2 user
	listen()

	// Start client connection

}

func listen() {
	listener, err := net.Listen(Protocol, NetWorkIp+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	var clients []net.Conn
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn, &clients)
		//mu.Lock()
		//if len(clients) < 2 {
		//	clients = append(clients, conn)
		//	fmt.Printf("Client %d connected: %s\n", len(clients), conn.RemoteAddr())
		//	go handleClient(conn, &clients)
		//} else {
		//	fmt.Println("Server already has two clients connected. Rejecting additional connections.")
		//	conn.Close()
		//}
		//mu.Unlock()
		//
		//if len(clients) == 2 {
		//	fmt.Println("Two clients connected. You can now start chatting!")
		//}
	}
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

func handleClient(conn net.Conn, clients *[]net.Conn) {
	for {
		//fmt.Fprintf(conn, "how many cards you want to send? ")
		//	scanner := bufio.NewScanner(conn)
		buffer := make([]byte, 1024)

		// Read data into the buffer
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}
		message := string(buffer[:bytesRead])
		cleaned := replacement(message)
		fmt.Printf("Received: %s", message)

		//if !scanner.Scan() {
		//	fmt.Println(" error reading from client")
		//	return
		//}
		//count, err := strconv.Atoi(scanner.Text())
		//if err != nil || count < 1 {
		//	fmt.Fprintf(conn, "invalid number: %v\n", err)
		//	return
		//}
		//var smsList []SmS
		//for t := 0; t < count; t++ {
		//	questions := []string{
		//		"From : ",
		//		"To : ",
		//		"Amount : ",
		//		"Bank :",
		//	}
		//	var sms SmS
		//	for i, question := range questions {
		//		fmt.Fprintf(conn, question+"\n")
		//		scanner := bufio.NewScanner(conn)
		//
		//		if scanner.Scan() {
		//			response := scanner.Text()
		//			switch i {
		//			case 0:
		//				sms.From = response
		//			case 1:
		//				sms.To = response
		//			case 2:
		//				amount, _ := strconv.Atoi(response)
		//				sms.Amount = amount
		//			case 3:
		//				sms.Bank = response
		//			}
		//			fmt.Println("Client said :", response)
		//		}
		//	}
		//	sms.OperationTime = time.Now()
		//	smsList = append(smsList, sms)
		var clientID string
		//response := cleaned
		mu.Lock()
		//for _, client := range *clients {
		//	if client == conn {
		clientID = conn.RemoteAddr().String()
		//messageWithID = fmt.Sprintf("Client %s says:\n", clientID)
		//for _, sms := range smsList {
		//	formattedSms := fmt.Sprintf("{#%s^%s$%d|%s-%s}\n", sms.From, sms.To, sms.Amount, sms.Bank, sms.OperationTime)
		//	messageWithID += formattedSms
		//}
		_, writeErr := conn.Write([]byte(clientID + " " + cleaned))
		if writeErr != nil {
			log.Fatal("Error writing to client:", clientID, writeErr)
		}

		fmt.Println(clientID + " " + cleaned)
		mu.Unlock()
	}
}

func removeClient(conn net.Conn, clients *[]net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	var updatedClients []net.Conn
	for _, client := range *clients {
		if client != conn {
			updatedClients = append(updatedClients, client)
		}
	}
	fmt.Println("Client disconnected:", conn.RemoteAddr())
	*clients = updatedClients
}
