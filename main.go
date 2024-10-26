package main

import (
	"bufio"
	"fmt"
	"net"
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

	// Continuously accept clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Only allow two clients to connect at any time
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

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", conn.RemoteAddr(), err)
			// Remove client from the list if it disconnects
			removeClient(conn, clients, mu)
			return
		}

		// Relay the message to the other client
		mu.Lock()
		for _, client := range *clients {
			if client != conn {
				_, writeErr := client.Write([]byte("payam az client digar : " + message))
				if writeErr != nil {
					fmt.Println("Error writing to client:", client.RemoteAddr(), writeErr)
				}
			}
		}
		mu.Unlock()
	}
}

// Remove client from the clients slice
func removeClient(conn net.Conn, clients *[]net.Conn, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	var updatedClients []net.Conn
	for _, client := range *clients {
		if client != conn {
			updatedClients = append(updatedClients, client)
		}
	}
	*clients = updatedClients
}
