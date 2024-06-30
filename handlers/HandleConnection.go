package handlers

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"sync"
    "time"
    "log"
)


const (
	MaxClients = 10
	DefaultPort = "8989"
    MaxLength = 20
    Usagemessage = "[USAGE]: ./TCPChat $port"
)

var (
	Clients = make(map[net.Conn]string)
	Messages =[]string{}
	ClientsMutex sync.Mutex
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
    //we cannot direclty convert slices of bytes to array of strings
    //loop over the array and convert each line (string) to slices of bytes
    length := len(WelcomeMessage())
    for i := 0; i < length; i++ {
        _, err := conn.Write([]byte(WelcomeMessage()[i] ))
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

    
	if len(name) > MaxLength{
		conn.Write([]byte(fmt.Sprintf("Name cannot exceed %d characters.\n", MaxLength)))
        goto name
	}

    for _, char := range name {
        if char < 'A' || char >  'Z' && char < 'a' || char > 'z'{
            conn.Write([]byte("Name cannot contain special characters or numbers.\n"))
			// handleConnection(conn) // Prompt user again
            goto name
        }
    }

	if name == "" {
		conn.Write([]byte("Name cannot be empty.\n"))
        goto name
	}

    if IsUsernameTaken(name) {
        conn.Write([]byte("Username is already taken, try another one.\n"))
        goto name
    }

	ClientsMutex.Lock()
	Clients[conn] = name
	ClientsMutex.Unlock()

	Broadcast(fmt.Sprintf("\n%s has joined our chat...", name), conn)

	SendPreviousMessages(conn)

	scanner := bufio.NewScanner(conn)

    go func() {
        conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), Clients[conn])))

        for scanner.Scan() {
            msg := scanner.Text()
            msg = strings.TrimSpace(msg)
            if msg == "" {
                conn.Write([]byte("Message cannot be empty.\n"))
                //Print timestamp and name again
                conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), Clients[conn])))
                
                continue
            }

                // check is message contains command to change name
                if strings.HasPrefix(msg, "/name ") {
                    // Handle /name command
                    newUsername := strings.TrimSpace(strings.TrimPrefix(msg, "/name "))
                    for {
                        if err := ValidateName(newUsername); err != nil {
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
                        if IsUsernameTaken(newUsername) {
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
                        oldName := Clients[conn]
                        Clients[conn] = newUsername
                        name = newUsername
        
                        // Broadcast username change
                        Broadcast(fmt.Sprintf("\n%s is now known as %s...", oldName, newUsername), conn)
                        //print timestamp and name again
                        conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), Clients[conn])))
                        break
                    }
                    continue // Continue handling messages
                }

            timestamp := time.Now().Format("2006-01-02 15:04:05")
            timestampedMessage := fmt.Sprintf("\r[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), Clients[conn], msg)
            ClientsMutex.Lock()
            Messages = append(Messages, timestampedMessage)
            ClientsMutex.Unlock()
            Broadcast(timestampedMessage, conn)
            conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", timestamp, Clients[conn])))

        }

        if err := scanner.Err(); err!= nil {
            log.Printf("Error reading from connection: %v\n", err)
        }

        ClientsMutex.Lock()
        delete(Clients, conn)
        ClientsMutex.Unlock()
        Broadcast(fmt.Sprintf("\n%s has left our chat...", name), conn)
    }()
    for {
        time.Sleep(1 *time.Second)
    }

}
