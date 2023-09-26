package main

import (
	"fmt"
	"net"
	"time"
)

func unicast(name, recipient, content string) {
	printHorizontalLine()
	fmt.Println(GreenColor + centerText("Sending Transmission", 40) + ResetColor)
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		fmt.Println("Error connecting to recipient:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Status: Connected")
	fmt.Println("Protocol: unicast")
	fmt.Printf("Recipient: %s\n", recipient)
	message := fmt.Sprintf("%s %s %s\n", Unicast, name, content)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	fmt.Println("Message Status: Sent")
	fmt.Println(GreenColor + centerText("End of Transmission", 40) + ResetColor)
	printHorizontalLine()
}

func multicast(name string, people map[string]string, content string, protocol string, timeout int, maxAttempts int) {
	message := fmt.Sprintf("%s %s %s\n", protocol, name, content)
	for _, recipient := range people {
		fmt.Println(GreenColor + centerText("Sending Transmission", 40) + ResetColor)
		fmt.Println("Status: Connected")
		fmt.Printf("Protocol: %s\n", protocol)
		fmt.Printf("Recipient: %s\n", recipient)
		for i := 0; i < maxAttempts; i++ {
			conn, err := net.Dial("tcp", recipient)
			if err != nil {
				fmt.Println("Error connecting to recipient:", err)
				return
			}
			defer conn.Close()
			_, err = conn.Write([]byte(message))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
			fmt.Println("Message Status: Sent - Waiting for Acknowledgement")
			timeoutACK := time.Duration(timeout)*time.Second + (time.Duration(i) * time.Second)
			fmt.Println("Timeout: ", timeoutACK)
			conn.SetReadDeadline(time.Now().Add(timeoutACK))
			buffer := make([]byte, 1024)
			_, err = conn.Read(buffer)
			if err != nil {
				fmt.Println("Timeout: No Acknowledgement Received, Retrying...")
			} else {
				fmt.Println("Message Status: Acknowledged")
				break
			}
		}
		fmt.Println(GreenColor + centerText("End of Transmission", 40) + ResetColor)
		printHorizontalLine()
	}
}
