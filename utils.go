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

func printStartScreen(nodes map[string]Node) {
	printHorizontalLine()
	fmt.Println(BlueColor + "Welcome " + nodes["thisNode"].name + ResetColor)
	fmt.Printf(YellowColor+"Your IP is: %v\n"+ResetColor, nodes["thisNode"].ip)
	fmt.Println(MagentaColor + "People found in the file: " + ResetColor)
	for _, node := range nodes {
		if !node.isThisNode {
			fmt.Printf("%s: %s\n", node.name, node.ip)
		}
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
	fmt.Println("Type '" + RedColor + "unicast" + ResetColor + "' to send a message to a specific person")
	fmt.Println("Type '" + RedColor + "clear" + ResetColor + "' to clear the screen")
	fmt.Println("Type '" + RedColor + "help" + ResetColor + "' to see this list of commands")
	fmt.Println("Type '" + RedColor + "get-delay" + ResetColor + "' to get the current delay")
	fmt.Println("Type '" + RedColor + "set-delay" + ResetColor + "' to set the current delay")
	fmt.Println("Type '" + RedColor + "status" + ResetColor + "' to get the nodes status")
	fmt.Println("Type '" + RedColor + "set-multicast-delay" + ResetColor + "' to set the multicast delay")
	fmt.Println("Type '" + RedColor + "set-unicast-delay" + ResetColor + "' to set the unicast delay")
	printHorizontalLine()
}
