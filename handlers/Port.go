package handlers

import "os"

func GetPort() string {
    if len(os.Args) < 2 {
        return DefaultPort
    }
    return os.Args[1]
}