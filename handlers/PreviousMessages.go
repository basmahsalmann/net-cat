package handlers

import (
	"net"
)

func SendPreviousMessages(conn net.Conn) {
    ClientsMutex.Lock()
    defer ClientsMutex.Unlock()
    for _, msg := range Messages {
        conn.Write([]byte(msg + "\n"))
    }
}