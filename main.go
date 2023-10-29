package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var delay int = 1
var multicastDelay = 0

type Node struct {
	name            string
	ip              string
	isAlive         bool
	isThisNode      bool
	expectedTimeout float64
	lastHeartbeat   time.Time
}

func NewNode(name string, ip string, isAlive bool, isThisNode bool, expectedTimeout float64) Node {
	return Node{name, ip, isAlive, isThisNode, expectedTimeout, time.Time{}}
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func getArgs() (string, string) {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide the arguments")
		os.Exit(1)
	}
	name := args[0]
	fileName := args[1]
	return name, fileName
}

func readFile(fileName string, thisName string) map[string]Node {
	nodes := make(map[string]Node)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			ip := parts[0]
			name := parts[1]
			if name == thisName {
				nodes["thisNode"] = NewNode(name, ip, true, true, DefaultIntervalForHeartbeat)
			} else {
				nodes[name] = NewNode(name, ip, true, false, DefaultIntervalForHeartbeat)
			}
		}
	}
	return nodes
}

func mainLoop(nodes map[string]Node) {
	for {
		protocol := readInput("> ")
		switch protocol {
		case "exit":
			os.Exit(0)
		case Multicast:
			fmt.Printf("[%s] Enter the message: ", protocol)
			message := readInput("")
			multicast(nodes, message)
		case Unicast:
			fmt.Println("Received unicast")
			recipientName := readInput("[unicast] Enter the name of the person: ")
			message := readInput("[unicast] Enter the message: ")
			unicast(recipientName, nodes, message)
		case "clear":
			fmt.Printf("\x1b[3;J\x1b[H\x1b[2J")
		case "help":
			printCommands()
		case "get-delay":
			fmt.Println(delay)
		case "status":
			for _, node := range nodes {
				if !node.isThisNode {
					fmt.Printf("\n%s: %s\n", node.name, node.ip)
					fmt.Printf("isAlive: %v\n", node.isAlive)
					fmt.Printf("expectedTimeout: %v seconds\n", node.expectedTimeout)
					fmt.Printf("lastHeartbeat: %v\n", node.lastHeartbeat)
				}
			}

		case "set-delay":
			fmt.Print("Enter the delay: ")
			delayStr := readInput("")
			delayInt, err := strconv.Atoi(delayStr)
			if err != nil {
				fmt.Println("Error converting delay to int:", err)
				return
			}
			delay = delayInt
		case "set-multicast-delay":
			fmt.Print("Enter the multicast delay: ")
			delayStr := readInput("")
			delayInt, err := strconv.Atoi(delayStr)
			if err != nil {
				fmt.Println("Error converting delay to int:", err)
				continue
			}
			multicastDelay = delayInt
		case "get-multicast-delay":
			fmt.Println(multicastDelay)

		default:
			fmt.Println("command not found")
		}
	}
}

func main() {
	thisName, fileName := getArgs()
	nodes := readFile(fileName, thisName)
	printStartScreen(nodes)
	go mainLoop(nodes)
	go heartbeat(nodes)
	// go checkHeartbeat(nodes)
	startServer(nodes)
}
