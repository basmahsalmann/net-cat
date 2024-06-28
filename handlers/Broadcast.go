package handlers

import (
	"fmt"
	"net"
	"time"
)



func Broadcast(message string, sender net.Conn) {
    ClientsMutex.Lock()
    defer ClientsMutex.Unlock()
    for conn := range Clients {
        // name := clients[conn]
        if conn != sender {
            conn.Write([]byte(message + "\n"))
            //when updating the client, we need to print the timestamp and name again
            conn.Write([]byte(fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), Clients[conn])))
        }
    }
}