package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// startServer starts the TCP server to receive incoming messages.
func startServer(ip string) {
	ln, err := net.Listen("tcp", ip)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
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

// handleConnection handles an incoming connection.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

	fmt.Printf("\n======receiving message======\n%s=====end of transmission=====\n", message)
}

// sendMessage sends a chat message to the specified recipient.
func sendMessage(name, recipient, content string) {
	fmt.Println("======sending message======")
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		fmt.Println("Error connecting to recipient:", err)
		return
	}
	defer conn.Close()
	conn.Write([]byte(name + ": " + content + "\n"))
	fmt.Printf("======message sent to %s======\n", recipient)

	fmt.Println("=====end of transmission=====")
}

// readInput reads a line of input from the user.
func readInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// function to deal with the parameters: it should expect two parameters, the name of the file and the name of the person
// it should return both the fileName and the name of the person
func getArgs() (string, string) {
	// read arguments from command line
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide an argument")
		return "", ""
	}
	fileName := args[0]
	name := args[1]
	return fileName, name
}

func readFile(fileName string) map[string]string {
	// read file
	lines := []string{}
	hosts, _ := os.Open(fileName)
	defer hosts.Close()
	// print file content
	scanner := bufio.NewScanner(hosts)
	// define a map to store people
	personMap := make(map[string]string)

	for scanner.Scan() {
		// append all lines to lines
		lines = append(lines, scanner.Text())
	}
	for _, line := range lines {
		// split line by space
		words := strings.Split(line, " ")
		// store in the map, name is the key and ip is the value
		personMap[words[1]] = words[0]
	}

	return personMap
}

func main() {
	name, fileName := getArgs()
	people := readFile(fileName)
	thisIP := people[name]
	delete(people, name)
	fmt.Println("this === " + thisIP)
	for _, ip := range people {
		fmt.Println(ip)
	}
	go startServer(thisIP)

	for {
		message := readInput("> ")
		for _, ip := range people {
			sendMessage(name, ip, message)
		}
	}
}
