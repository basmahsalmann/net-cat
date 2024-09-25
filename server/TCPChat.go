package main

import (
	"fmt"
	"log"
	"net"
	"os"

	n "topchat/handlers"
)

func main() {
	port := n.GetPort()
	//this gets the external ip address of the machine and prints it to make it easier to connect to the server
	//but we can go extra steps and use "ip a" in terminal to get the ip
	ip := n.GetLocalIP()

      //Usage: ./TCPChat $port
	if len(os.Args) > 2 {
        fmt.Println(n.Usagemessage)
        os.Exit(1)
    }

	// create server to listen for connections
	//"0.0.0.0" will allow incoming connections from any network interface, local or external
	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error setting up listener: %v\n", err)
	}
	defer ln.Close()

	fmt.Printf("Listening on port %s:%s\n", ip, port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		n.ClientsMutex.Lock()
		if len(n.Clients) >= n.MaxClients {
			n.ClientsMutex.Unlock()
			conn.Write([]byte("Server is full. Try again later.\n"))
			conn.Close()
			continue
		}
		n.ClientsMutex.Unlock()

		// handle each client
		go n.HandleConnection(conn)
	}
}
