package main

import (
	"fmt"
	"net"
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

func multicast(name string, people map[string]string, message string) {
	// implement a multicast protocol this will consist of sending a message to all the people in the map
	for _, ip := range people {
		unicast(name, ip, message)
	}
}
