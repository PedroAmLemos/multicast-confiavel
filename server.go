package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func startServer(ip string) {
	const maxRetries = 3
	retries := 0

	for retries < maxRetries {
		ln, err := net.Listen("tcp", ip)
		if err != nil {
			fmt.Printf("Error starting server: %s. Attempt %d/%d\n", err, retries+1, maxRetries)
			retries++
			time.Sleep(2 * time.Second)
			continue
		}

		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go handleConnection(conn)
		}
	}

	if retries == maxRetries {
		fmt.Println("Failed to start the server after multiple attempts. Exiting.")
		os.Exit(1)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

	parts := strings.SplitN(message, " ", 2)
	if len(parts) < 2 {
		fmt.Println("Received malformed message:", message)
		return
	}
	protocol := parts[0]
	actualMessage := parts[1]

	fmt.Println()
	printHorizontalLine()
	fmt.Println(BlueColor + centerText("Received Message", 40) + ResetColor)
	fmt.Printf("Protocol: %s\n", protocol)
	fmt.Println("Message Content:")
	fmt.Printf("%s", actualMessage)
	fmt.Println(BlueColor + centerText("End of Message", 40) + ResetColor)
	printHorizontalLine()
	fmt.Print("> ")
}
