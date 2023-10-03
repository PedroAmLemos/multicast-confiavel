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
	// make the commands with color
	fmt.Println("Type '" + RedColor + "exit" + ResetColor + "' to exit")
	fmt.Println("Type '" + RedColor + "list" + ResetColor + "' to list all people")
	fmt.Println("Type '" + RedColor + "multicast" + ResetColor + "' to send a message to all people")
	fmt.Println("Type '" + RedColor + "multicast-delay" + ResetColor + "' to send a message to all people and use the delayed response handler")
	fmt.Println("Type '" + RedColor + "unicast" + ResetColor + "' to send a message to a specific person")
	fmt.Println("Type '" + RedColor + "clear" + ResetColor + "' to clear the screen")
	fmt.Println("Type '" + RedColor + "help" + ResetColor + "' to see this list of commands")
	fmt.Println("Type '" + RedColor + "get-delay" + ResetColor + "' to get the current delay")
	fmt.Println("Type '" + RedColor + "set-delay" + ResetColor + "' to set the current delay")
	printHorizontalLine()
}
