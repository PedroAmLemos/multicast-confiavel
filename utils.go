package main

import (
	"fmt"
	"strings"
)

func centerText(text string, width int) string {
	padding := width - len(text)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding
	return strings.Repeat("=", leftPadding) + " " + text + " " + strings.Repeat("=", rightPadding)
}

func printHorizontalLine() {
	fmt.Println(CyanColor + "------------------------------------------------" + ResetColor)
}

func printStartScreen(name string, thisIP string, people map[string]string) {
	printHorizontalLine()
	fmt.Println(BlueColor + "Welcome " + name + ResetColor)
	fmt.Printf(YellowColor+"Your IP is: %s\n"+ResetColor, thisIP)
	fmt.Println(MagentaColor + "People found in the file: " + ResetColor)
	for name, ip := range people {
		fmt.Printf("%s: %s\n", name, ip)
	}
	fmt.Println("Type " + RedColor + "'help'" + ResetColor + " to see the list of commands")
	printHorizontalLine()
}

func printCommands() {
	printHorizontalLine()
	fmt.Println("Type 'exit' to exit")
	fmt.Println("Type 'list' to list all people")
	fmt.Println("Type 'multicast' to send a message to all people")
	fmt.Println("Type 'multicast-delay' to send a message to all people and use the delayed response handler")
	fmt.Println("Type 'unicast' to send a message to a specific person")
	fmt.Println("Type 'clear' to clear the screen")
	fmt.Println("Type 'help' to see this list of commands")
	fmt.Println("Type 'get-delay' to get the current delay")
	fmt.Println("Type 'set-delay' to set the current delay")
	printHorizontalLine()
}
