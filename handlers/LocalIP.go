package handlers

import (
	"net"
	"os"
	"log"
)



func GetLocalIP() string {
    var localAddress string
    //if the user has specified a local address, use that (should be local & valid or an error will be displayed)
    if len(os.Args) > 2 {
        localAddress = os.Args[2]
        return localAddress
    }
    //otherwise, use the local address of the machine running the program
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddress = (conn.LocalAddr().(*net.UDPAddr)).IP.String()

    return localAddress
}