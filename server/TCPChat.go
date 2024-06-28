package main

import (
	"fmt"
	"log"
	"net"
    n "topchat/handlers"
)

func main() {
    port := n.GetPort()
    ip := n.GetLocalIP()

    //create server to listen for connections
    ln, err := net.Listen("tcp", ip+ ":"+port)
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

        //handle each client
        go n.HandleConnection(conn)
    }
}

