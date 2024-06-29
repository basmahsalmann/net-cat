package handlers

import (
	"net"
	"log"
    "os"
    "fmt"
)



func GetLocalIP() string {
    var localAddress string

    //Usage: ./TCPChat $port
    if len(os.Args) > 2 {
        fmt.Println(Usagemessage)
        os.Exit(1)
    }
    //use the local address of the machine running the program
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddress = (conn.LocalAddr().(*net.UDPAddr)).IP.String()

    return localAddress
}