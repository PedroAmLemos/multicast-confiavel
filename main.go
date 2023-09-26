package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func startServer(ip string) {
	ln, err := net.Listen("tcp", ip)
	if err != nil {
		panic(err)
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

func handleConnection(conn net.Conn) {
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

	fmt.Printf("\n======receiving message======\n%s=====end of transmission=====\n> ", message)
}

func sendMessage(name, recipient, content string, protocol string) {
	fmt.Println("======sending transmission======")
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		fmt.Println("Error connecting to recipient:", err)
		return
	}
	defer conn.Close()
	fmt.Println("...connected...")
	fmt.Printf("...protocol: %s...\n", protocol)
	switch protocol {
	case "unicast":
		fmt.Printf("...recipient: %s...\n", recipient)
	case "multicast":
		fmt.Println("...recipient: all...")
	}
	message := fmt.Sprintf("%s %s: %s\n", protocol, name, content)
	conn.Write([]byte(message))
	fmt.Println("...sent...")
	fmt.Println("=====end of transmission=====")
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
	fileName := args[0]
	name := args[1]
	return fileName, name
}

func readFile(fileName string) map[string]string {
	lines := []string{}
	hosts, _ := os.Open(fileName)
	defer hosts.Close()
	scanner := bufio.NewScanner(hosts)
	personMap := make(map[string]string)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for _, line := range lines {
		words := strings.Split(line, " ")
		personMap[words[1]] = words[0]
	}

	return personMap
}

func multicast(name string, people map[string]string, message string) {
}

func printStartScreen(name string, thisIP string, people map[string]string) {
	fmt.Println("Welcome " + name)
	fmt.Printf("Your IP is: %s\n", thisIP)
	fmt.Println("People found in the file: ")
	for name, ip := range people {
		fmt.Printf("%s: %s\n", name, ip)
	}
	fmt.Println("Type 'commands' to see the list of commands")
}

func printCommands() {
	fmt.Println("Type 'exit' to exit")
	fmt.Println("Type 'list' to list all people")
	fmt.Println("Type 'multicast' to send a message to all people")
	fmt.Println("Type 'unicast' to send a message to a specific person")
	fmt.Println("Type 'clear' to clear the screen")
}

func main() {
	name, fileName := getArgs()
	people := readFile(fileName)
	thisIP := people[name]
	delete(people, name)
	printStartScreen(name, thisIP, people)
	go func() {
		for {
			command := readInput("> ")
			switch command {
			case "exit":
				os.Exit(0)
			case "list":
				for name, ip := range people {
					fmt.Printf("%s: %s\n", name, ip)
				}
			case "multicast":
				message := readInput("[multicast] Enter the message: ")
				multicast(name, people, message)
			case "unicast":
				recipient := readInput("[unicast] Enter the name of the person: ")
				message := readInput("[unicast] Enter the message: ")
				recipientIP := people[recipient]
				sendMessage(name, recipientIP, message, command)
			case "clear":
				os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
			case "commands":
				printCommands()
			default:
				fmt.Println("command not found")
			}
		}
	}()
	startServer(thisIP)
}
