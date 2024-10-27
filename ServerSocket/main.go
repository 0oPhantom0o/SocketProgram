package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

const NetWorkIp = "0.0.0.0:8181"
const Protocol = "tcp"

func main() {
	listener, err := net.Listen(Protocol, NetWorkIp)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8181...")
	var clients []net.Conn
	var mu sync.Mutex

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		mu.Lock()
		if len(clients) < 2 {
			clients = append(clients, conn)
			fmt.Printf("Client %d connected: %s\n", len(clients), conn.RemoteAddr())
			go handleClient(conn, &clients, &mu)
		} else {
			fmt.Println("Server already has two clients connected. Rejecting additional connections.")
			conn.Close()
		}
		mu.Unlock()

		if len(clients) == 2 {
			fmt.Println("Two clients connected. You can now start chatting!")
		}
	}
}

func handleClient(conn net.Conn, clients *[]net.Conn, mu *sync.Mutex) {
	defer conn.Close()

	for {
		message, err := convertOutput(conn)
		if err != nil {
			fmt.Println("Error reading from client:", conn.RemoteAddr(), err)
			removeClient(conn, clients, mu)
			return
		}

		mu.Lock()
		for _, client := range *clients {
			if client != conn {

				clientID := conn.RemoteAddr().String()
				messageWithID := fmt.Sprintf("Client %s says:\n%s", clientID, message)
				//	fmt.Printf(" message to client %s:\n%s", client.RemoteAddr(), messageWithID)
				_, writeErr := client.Write([]byte(messageWithID))
				if writeErr != nil {
					fmt.Println("Error writing to client:", client.RemoteAddr(), writeErr)
				}
			}
		}
		mu.Unlock()
	}
}

func convertOutput(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	//this line is for find end of message  . example : if client press enter data will move to server . if you want to ask for 10 line message you must make a loop here
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	formattedMessage := formatNames(data)
	return formattedMessage, nil
}
func formatNames(input string) string {
	// in this function we will send 3 names in 1 line but it will make it to 3 line string
	// Remove any extra whitespace at the start and end of the input
	input = strings.TrimSpace(input)

	// Split the input string by spaces
	names := strings.Fields(input)

	// Join each name with a newline character
	return strings.Join(names, "\n")
}
func removeClient(conn net.Conn, clients *[]net.Conn, mu *sync.Mutex) {
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
