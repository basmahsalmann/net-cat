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
    maxLength = 20
    // ip = "10.1.201.169"
)

var (
	clients = make(map[net.Conn]string)
	messages =[]string{}
	clientsMutex sync.Mutex
)

func main() {
    port := getPort()
    ip := GetLocalIP().String()
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

name:
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	name, err := bufio.NewReader(conn).ReadString('\n')
	if err!= nil {
		// log.Printf("Error reading name: %v\n", err)
		return
	}

	name = strings.TrimSpace(name)

    
	if len(name) > maxLength{
		conn.Write([]byte(fmt.Sprintf("Name cannot exceed %d characters.\n", maxLength)))
        goto name
		return
	}

    for _, char := range name {
        if char < 'A' || char >  'Z' && char < 'a' || char > 'z'{
            conn.Write([]byte("Name cannot contain special characters or numbers.\n"))
			// handleConnection(conn) // Prompt user again
            goto name
			return
        }
    }

	if name == "" {
		conn.Write([]byte("Name cannot be empty.\n"))
        goto name
        return
	}

    if isUsernameTaken(name) {
        conn.Write([]byte("Username is already taken, try another one.\n"))
        goto name
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
            msg = strings.TrimSpace(msg)
            if msg == "" {
                conn.Write([]byte("Message cannot be empty.\n"))
                //Print timestamp and name again
                conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), name)))
                
                continue
            }

                // check is message contains command to change name
                if strings.HasPrefix(msg, "/name ") {
                    // Handle /name command
                    newUsername := strings.TrimSpace(strings.TrimPrefix(msg, "/name "))
                    for {
                        if err := validateName(newUsername); err != nil {
                            conn.Write([]byte(fmt.Sprintf("%s\n", err)))
                            conn.Write([]byte("[ENTER ANOTHER NAME]: "))
                            input, err := bufio.NewReader(conn).ReadString('\n')
                            if err != nil {
                                log.Printf("Error reading name: %v\n", err)
                                continue
                            }
                            newUsername = strings.TrimSpace(input)
                            continue
                        }
        
                        // Check if username is already taken
                        if isUsernameTaken(newUsername) {
                            conn.Write([]byte(fmt.Sprintf("Username %s is already taken. Please choose another.\n", newUsername)))
                            conn.Write([]byte("[ENTER ANOTHER NAME]: "))
                            input, err := bufio.NewReader(conn).ReadString('\n')
                            if err != nil {
                                log.Printf("Error reading name: %v\n", err)
                                continue
                            }
                            newUsername = strings.TrimSpace(input)
                            continue
                        }
        
                        // Update username
                        oldName := clients[conn]
                        clients[conn] = newUsername
        
                        // Broadcast username change
                        broadcast(fmt.Sprintf("\n%s is now known as %s...", oldName, newUsername), conn)
                        break
                    }
                    continue // Continue handling messages
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
    }()
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

func GetLocalIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddress := conn.LocalAddr().(*net.UDPAddr)

    return localAddress.IP
}

func isUsernameTaken(username string) bool {
    clientsMutex.Lock()
    defer clientsMutex.Unlock()
    for _, existingName := range clients {
        if existingName == username {
            return true
        }
    }
    return false
}

func validateUsername(name string) error {
    if len(name) > maxLength {
        return fmt.Errorf("Name cannot exceed %d characters", maxLength)
    }
    for _, char := range name {
        if (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
            return fmt.Errorf("Name cannot contain special characters or numbers")
        }
    }
    return nil
}

func validateName(name string) error {
	if len(name) > maxLength {
		return fmt.Errorf("Name cannot exceed %d characters", maxLength)
	}
	for _, char := range name {
		if (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
			return fmt.Errorf("Name cannot contain special characters or numbers")
		}
	}
	return nil
}