package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func startServer(nodes map[string]Node) {
	ip := nodes["thisNode"].ip
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
			go handleConnection(conn, nodes)
		}
	}

	if retries == maxRetries {
		fmt.Println("Failed to start the server after multiple attempts. Exiting.")
		os.Exit(1)
	}
}

func handleMulticast(name string, message string, protocol string, conn net.Conn) {
	time.Sleep(time.Second * time.Duration(multicastDelay))
	fmt.Println()
	printHorizontalLine()
	fmt.Println(BlueColor + centerText("Received Message", 40) + ResetColor)
	fmt.Printf("Protocol: %s\n", protocol)
	fmt.Printf("Message Content: %s", message)
	fmt.Printf("Sender: %s\n", name)
	if protocol == MulticastDelay {
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

func handleConnection(conn net.Conn, nodes map[string]Node) {
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
	case Heartbeat:
		node := nodes[name]
		heartbeatTime := time.Now()
		node.isAlive = true
		if node.lastHeartbeat.IsZero() {
			fmt.Printf("\nFirst heartbeat from %v... Marking it as alive\n> ", node.name)
		} else {
			expected := node.lastHeartbeat.Add(time.Duration(node.expectedTimeout) * time.Second)
			if heartbeatTime.After(expected) {
				diff := heartbeatTime.Sub(expected)
				if diff > time.Duration(time.Millisecond*2) {
					node.expectedTimeout = node.expectedTimeout + diff.Seconds()
					fmt.Printf("\nReceived heartbeat from %v after expected timeout. Adjusting timeout to %v.\n> ", node.name, node.expectedTimeout)
				}
			} else if heartbeatTime.Before(expected) {
				diff := expected.Sub(heartbeatTime)
				if diff > time.Duration(time.Millisecond*2) {
					node.expectedTimeout = node.expectedTimeout - diff.Seconds()
					fmt.Printf("\nReceived heartbeat from %v before expected timeout. Adjusting timeout to %v.\n> ", node.name, node.expectedTimeout)
				}
			}
		}
		node.lastHeartbeat = heartbeatTime
		nodes[name] = node

	case Multicast:
		handleMulticast(name, actualMessage, protocol, conn)
	case Unicast:
		handleUnicast(name, actualMessage)

	default:
		fmt.Println("Received unknown protocol: ", protocol)
	}
}

func heartbeat(nodes map[string]Node) {
	thisName := nodes["thisNode"].name
	count := 0
	for {
		count++
		for _, node := range nodes {
			if !node.isThisNode && (node.isAlive || count%10 == 0) {
				conn, err := net.Dial("tcp", node.ip)
				if err != nil {
					fmt.Printf("\nError connecting to %v: %v\n> ", node.name, err)
					if node.isAlive {
						fmt.Printf("\nMarking %v as dead and continuing...\n> ", node.name)
					}
					currentNode := nodes[node.name]
					currentNode.isAlive = false
					currentNode.lastHeartbeat = time.Time{}
					currentNode.expectedTimeout = DefaultIntervalForHeartbeat
					nodes[node.name] = currentNode
					continue
				}
				hbMessage := fmt.Sprintf("%s %s PING\n", Heartbeat, thisName)
				_, err = conn.Write([]byte(hbMessage))
				if err != nil {
					fmt.Printf("\nError sending message: %v\n> ", err)
					continue
				}
				conn.Close()
			}
		}
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
		} else {
			time.Sleep(time.Duration(DefaultIntervalForHeartbeat) * time.Second)
		}
	}
}

func checkHeartbeat(nodes map[string]Node) {
	for {
		fmt.Printf("\nChecking heartbeat...\n> ")
		for _, node := range nodes {
			if !node.isThisNode && node.isAlive && node.lastHeartbeat != (time.Time{}) {
				expected := node.lastHeartbeat.Add(time.Duration(node.expectedTimeout*2) * time.Second)
				if time.Now().After(expected) {
					fmt.Printf("\nNode %v did not send a heartbeat for more then 2 * expectedInterval. Marking as dead.\n> ", node.name)
					currentNode := nodes[node.name]
					currentNode.isAlive = false
					nodes[node.name] = currentNode
				}
			}
		}
		time.Sleep(15 * time.Second)
	}
}
