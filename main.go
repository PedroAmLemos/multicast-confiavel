package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var delay int = 0

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

func readFile(fileName string) map[string]string {
	var lines []string
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
		protocol := readInput("> ")
		switch protocol {
		case "exit":
			os.Exit(0)
		case "list":
			for name, ip := range people {
				fmt.Printf("%s: %s\n", name, ip)
			}
		case Multicast, MulticastDelay:
			fmt.Printf("[%s] Enter the message: ", protocol)
			message := readInput("")
			fmt.Printf("[%s] Enter the max timeout (in seconds): ", protocol)
			timeout := readInput("")
			timeout = strings.TrimSpace(timeout)
			timeoutInt, err := strconv.Atoi(timeout)
			if err != nil {
				fmt.Println("Error converting timeout to int:", err)
				return
			}
			fmt.Printf("[%s] Enter the max attempts: ", protocol)
			attempts := readInput("")
			attempts = strings.TrimSpace(attempts)
			attemptsInt, err := strconv.Atoi(attempts)
			if err != nil {
				fmt.Println("Error converting attempts to int:", err)
				return
			}
			multicast(name, people, message, protocol, timeoutInt, attemptsInt)
		case Unicast:
			recipient := readInput("[unicast] Enter the name of the person: ")
			message := readInput("[unicast] Enter the message: ")
			recipientIP := people[recipient]
			unicast(name, recipientIP, message)
		case "clear":
			fmt.Printf("\x1b[3;J\x1b[H\x1b[2J")
		case "help":
			printCommands()

		case "get-delay":
			fmt.Println(delay)
		case "set-delay":
			fmt.Print("Enter the delay: ")
			delayStr := readInput("")
			delayInt, err := strconv.Atoi(delayStr)
			if err != nil {
				fmt.Println("Error converting delay to int:", err)
				return
			}
			delay = delayInt

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
