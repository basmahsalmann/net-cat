package main

import (
	"os"
	"fmt"
	"bufio"
	"sync"
	"time"
	"log"
	"net"
    "strings"
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
    //we cannot direclty convert slices of bytes to array of strings
    //loop over the array and convert each line (string) to slices of bytes
    length := len(welcomeMessage())
    for i := 0; i < length; i++ {
        _, err := conn.Write([]byte(welcomeMessage()[i] ))
        if err != nil {
            fmt.Println("Error sending message to client.")
            return
        }

    }
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	name, err := bufio.NewReader(conn).ReadString('\n')
	if err!= nil {
		log.Printf("Error reading name: %v\n", err)
		return
	}

	name = strings.TrimSpace(name)

	if name == "" {
		conn.Write([]byte("Name cannot be empty.\n"))
		return
	}

	clientsMutex.Lock()
	clients[conn] = name
	clientsMutex.Unlock()

	broadcast(fmt.Sprintf("\n%s has joined our chat...", name), conn)

	sendPreviousMessages(conn)

	scanner := bufio.NewScanner(conn)

    go func() {
        conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), name)))

        for scanner.Scan() {
            msg := scanner.Text()
            if msg == "" {
                conn.Write([]byte("Message cannot be empty.\n"))
                //Print timestamp and name again
                conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), name)))

                continue
            }
            timestamp := time.Now().Format("2006-01-02 15:04:05")
            timestampedMessage := fmt.Sprintf("\r[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, msg)
            clientsMutex.Lock()
            messages = append(messages, timestampedMessage)
            clientsMutex.Unlock()
            broadcast(timestampedMessage, conn)
            conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", timestamp, name)))

        }

        if err := scanner.Err(); err!= nil {
            log.Printf("Error reading from connection: %v\n", err)
        }

        clientsMutex.Lock()
        delete(clients, conn)
        clientsMutex.Unlock()
        broadcast(fmt.Sprintf("\n%s has left our chat...", name), conn)
    } ()
    for {
        time.Sleep(1 *time.Second)
        // conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), name)))
    }

}

func welcomeMessage() []string {
    //Create empty array of strings to store linux logo
	fileText := make([]string, 0)//
	file, err := os.Open("logo.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	//Reads from file
	scanner := bufio.NewScanner(file)

	//Return true if there is a line to read
	for scanner.Scan() {
		fileText = append(fileText, scanner.Text() + "\n") // appends line to fileText
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erorr reading file: %s", err)
	}

	return fileText
}

func broadcast(message string, sender net.Conn) {
    clientsMutex.Lock()
    defer clientsMutex.Unlock()
    for conn := range clients {
        name := clients[conn]
        if conn != sender {
            conn.Write([]byte(message + "\n"))
            //when updating the client, we need to print the timestamp and name again
            conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), name)))
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
