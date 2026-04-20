package main

import (
	"fmt"
	"os"

	"adbr.xx/anomchat/client"
	"adbr.xx/anomchat/server"
	"adbr.xx/anomchat/utils"
)

func main() {
	switch os.Args[1] {
	case "c":
		if len(os.Args) < 4 {
			fmt.Println("Please provide a host and a key")
			return
		}

		host := os.Args[2]
		key := os.Args[3]
		chat := utils.ChatInfo{Host: host, Key: key}
		client.InitializeClient(chat)
	case "s":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a port number")
			return
		}
		port := os.Args[2]

		pwd, err := utils.GenerateRandomPassword()
		if err != nil {
			fmt.Println("Error generating password:", err)
			return
		}
		server.InitializeTCPServer(port, pwd)
	default:
		fmt.Println("Invalid option. Use 'c' for client or 's' for server.")
	}
}
