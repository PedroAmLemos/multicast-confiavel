package main

import (
	"fmt"
	"net"
)

func sendMessage(name, recipient, content string, protocol string) {
	printHorizontalLine()
	fmt.Println(GreenColor + centerText("Sending Transmission", 40) + ResetColor)
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		fmt.Println("Error connecting to recipient:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Status: Connected")
	fmt.Printf("Protocol: %s\n", protocol)
	switch protocol {
	case Unicast:
		fmt.Printf("Recipient: %s\n", recipient)
	case Multicast:
		fmt.Println("Recipient: All")
	}
	message := fmt.Sprintf("%s %s: %s\n", protocol, name, content)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	fmt.Println("Message Status: Sent")
	fmt.Println(GreenColor + centerText("End of Transmission", 40) + ResetColor)
	printHorizontalLine()
}

func multicast(name string, people map[string]string, message string) {
}
