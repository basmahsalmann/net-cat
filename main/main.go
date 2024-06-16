package main

import (
	"os"
	"fmt"
	"bufio"
	"sync"
	"time"
	"log"
	"net"
)

const (
	maxClients = 10
	defaultPort = "8989"
)

var (
	clients = make(map[net.Conn]string)
	messages =[]string{}
	clientsMutex sync.Mutex
)

func main() {
    port := getPort()
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("Error setting up listener: %v\n", err)
    }
    defer ln.Close()

    fmt.Printf("Listening on port :%s\n", port)

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %v\n", err)
            continue
        }

        clientsMutex.Lock()
        if len(clients) >= maxClients {
            clientsMutex.Unlock()
            conn.Write([]byte("Server is full. Try again later.\n"))
            conn.Close()
            continue
        }
        clientsMutex.Unlock()

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte(welcomeMessage()))
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	name, err := bufio.NewReader(conn).ReadString('\n')
	if err!= nil {
		log.Printf("Error reading name: %v\n", err)
		return
	}

	name = name[:len(name)-1]

	if name == "" {
		conn.Write([]byte("Name cannot be empty.\n"))
		return
	}

	clientsMutex.Lock()
	clients[conn] = name
	clientsMutex.Unlock()

	broadcast(fmt.Sprintf("%s has joined our chat...\n", name), conn)

	sendPreviousMessages(conn)

	scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        msg := scanner.Text()
        if msg == "" {
            conn.Write([]byte("Message cannot be empty.\n"))
            continue
        }
        timestampedMessage := fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, msg)
        clientsMutex.Lock()
        messages = append(messages, timestampedMessage)
        clientsMutex.Unlock()
        broadcast(timestampedMessage, conn)
    }

    if err := scanner.Err(); err!= nil {
        log.Printf("Error reading from connection: %v\n", err)
    }

    clientsMutex.Lock()
    delete(clients, conn)
    clientsMutex.Unlock()
    broadcast(fmt.Sprintf("%s has left our chat...\n", name), conn)

}

func welcomeMessage() string {
    return `Welcome to TCP-Chat!`
}

func broadcast(message string, sender net.Conn) {
    clientsMutex.Lock()
    defer clientsMutex.Unlock()
    for conn := range clients {
        if conn != sender {
            conn.Write([]byte(message + "\n"))
        }
    }
}

func sendPreviousMessages(conn net.Conn) {
    clientsMutex.Lock()
    defer clientsMutex.Unlock()
    for _, msg := range messages {
        conn.Write([]byte(msg + "\n"))
    }
}

func getPort() string {
    if len(os.Args) < 2 {
        return defaultPort
    }
    return os.Args[1]
}
