package main

import (
	"fmt"
	"net"
	"time"
)

func unicast(receiverName string, nodes map[string]Node, content string) {
	receiverIp := nodes[receiverName].ip
	thisName := nodes["thisNode"].name
	printHorizontalLine()
	fmt.Println(GreenColor + centerText("Sending Transmission", 40) + ResetColor)
	conn, err := net.Dial("tcp", receiverIp)
	if err != nil {
		fmt.Printf("\nError connecting to %v: %v\n", receiverName, err)
		return
	}
	defer conn.Close()
	fmt.Println("Status: Connected")
	fmt.Println("Protocol: unicast")
	fmt.Printf("Recipient: %s - %s\n", receiverIp, receiverName)
	message := fmt.Sprintf("%s %s %s\n", Unicast, thisName, content)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message: ", err)
		return
	}
	fmt.Println("Message Status: Sent")
	fmt.Println(GreenColor + centerText("End of Transmission", 40) + ResetColor)
	printHorizontalLine()
}

func multicast(nodes map[string]Node, messageContent string) {
	message := fmt.Sprintf("%s %s %s\n", Multicast, nodes["thisNode"].name, messageContent)
	fmt.Println(GreenColor + centerText("Starting Multicast", 40) + ResetColor)
	fmt.Printf("\nMessage Content: %s\n", message)
	receivedACK := make(map[string]bool)
	for _, node := range nodes {
		printHorizontalLine()
		if node.isThisNode {
			continue
		}
		if node.isAlive {
			conn, err := net.Dial("tcp", node.ip)
			if err != nil {
				fmt.Printf("\nError connecting to %v: %v\n", node.name, err)
				continue
			}
			fmt.Printf("Status: Connected to %v - %v\nProtocol: Multicast\n", node.name, node.ip)
			_, err = conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("\nError sending message to %v: %v\n", node.name, err)
				continue
			}
			fmt.Printf("Message Status: Sent to %v\n", node.name)
			err = conn.SetReadDeadline(time.Now().Add(time.Duration(node.expectedTimeout) * time.Second * 2).Add(time.Millisecond * 2))
			if err != nil {
				fmt.Printf("\nError setting read deadline: %v\n", err)
				continue
			} else {
				fmt.Printf("Read deadline set to %v seconds (expected timeout (%v) * 2)\n", node.expectedTimeout*2, node.expectedTimeout)
				_, err = conn.Read(make([]byte, 1024))
				if err != nil {
					fmt.Printf("\nError reading from %v: %v\n", node.name, err)
					receivedACK[node.name] = false
					continue
				}
				fmt.Printf("Message Status: Received ACK from %v\n", node.name)
				receivedACK[node.name] = true
			}
		} else {
			fmt.Printf("Status: %v is not alive\nProtocol: Multicast\n", node.name)
			fmt.Printf("Message Status: Not sent to %v\n", node.name)
			fmt.Printf("Reason: %v is not alive\n", node.name)
			fmt.Printf("Continuing to next node...\n")
			receivedACK[node.name] = false
			continue
		}

	}
	fmt.Println(GreenColor + centerText("End of Multicast", 40) + ResetColor)
	printHorizontalLine()
	fmt.Printf("MULTICAST REPORT\n")
	for name, ack := range receivedACK {
		if ack {
			fmt.Printf("%v: ACK\n", name)
		} else {
			fmt.Printf("%v: NACK\n", name)
		}
	}
	printHorizontalLine()
}
