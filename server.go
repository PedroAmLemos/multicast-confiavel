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

func handleMulticast(name string, message string, protocol string, conn net.Conn, delay int) {
	fmt.Println()
	printHorizontalLine()
	fmt.Println(BlueColor + centerText("Received Message", 40) + ResetColor)
	fmt.Printf("Protocol: %s\n", protocol)
	fmt.Printf("Message Content: %s", message)
	fmt.Printf("Sender: %s\n", name)
	if delay > 0 {
		fmt.Printf("Starting delay: %d seconds\n", delay)
		delayTime := time.Now().Add(time.Duration(delay) * time.Second)
		time.Sleep(time.Until(delayTime))
		fmt.Printf("Ending delay: %d seconds\n", delay)
	}
	fmt.Println("Sending ACK back to", name)
	conn.Write([]byte("ACK\n"))
	fmt.Println(BlueColor + centerText("End of Message", 40) + ResetColor)
	printHorizontalLine()
	fmt.Print("> ")
}

func handleMulticastDelay(name string, message string) {
	fmt.Println("Received multicast delay message")
}

func handleUnicast(name string, message string) {
	fmt.Println()
	printHorizontalLine()
	fmt.Println(BlueColor + centerText("Received Message", 40) + ResetColor)
	fmt.Printf("Protocol: %s\n", Unicast)
	fmt.Printf("Message Content: %s", message)
	fmt.Printf("Sender: %s\n", name)
	fmt.Println(BlueColor + centerText("End of Message", 40) + ResetColor)
	printHorizontalLine()
	fmt.Print("> ")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

	parts := strings.SplitN(message, " ", 3)
	if len(parts) < 3 {
		fmt.Println("Received malformed message:", message)
		return
	}
	protocol := parts[0]
	name := parts[1]
	actualMessage := parts[2]
	switch protocol {
	case Multicast:
		handleMulticast(name, actualMessage, protocol, conn, 0)
	case MulticastDelay:
		handleMulticast(name, actualMessage, protocol, conn, 2)
	case Unicast:
		handleUnicast(name, actualMessage)

	default:
		fmt.Println("Received unknown protocol: ", protocol)
	}
}
