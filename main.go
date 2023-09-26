package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func mainLoop(people map[string]string, name string) {
	for {
		command := readInput("> ")
		switch command {
		case "exit":
			os.Exit(0)
		case "list":
			for name, ip := range people {
				fmt.Printf("%s: %s\n", name, ip)
			}
		case Multicast:
			message := readInput("[multicast] Enter the message: ")
			multicast(name, people, message)
		case Unicast:
			recipient := readInput("[unicast] Enter the name of the person: ")
			message := readInput("[unicast] Enter the message: ")
			recipientIP := people[recipient]
			unicast(name, recipientIP, message)
		case "clear":
			os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
		case "commands":
			printCommands()
		case "multicast delay":
			message := readInput("[multicast delay] Enter the message: ")
			multicast(name, people, message)
		default:
			fmt.Println("command not found")
		}
	}
}

func main() {
	name, fileName := getArgs()
	people := readFile(fileName)
	thisIP := people[name]
	delete(people, name)
	printStartScreen(name, thisIP, people)
	go mainLoop(people, name)
	startServer(thisIP)
}
