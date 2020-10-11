package main

import (
	"bufio"
	"fmt"
	//"log"
	"net"
	"os"
	//"strconv"
	"strings"
)

var connmap = make(map[string]net.Conn)

func createUser(c net.Conn, q string) string {
	c.Write([]byte(string(q)))
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return netData
		}

		username := strings.TrimSpace(string(netData))

		if username == "" {
			c.Write([]byte(string("You need to enter a username: ")))
		} else {
			c.Write([]byte(string("Welcome " + username + "\n")))
			return username
		}
	}
	c.Close()
	return string("testuser")
}

func enterRoom(c net.Conn, username string) {
	// tell the user who else is in the room
	for key, _ := range connmap {
		c.Write([]byte(string("\n" + key + " is in the room\n# ")))
	}
	// the the other users who has joined the room
	for _, value := range connmap {
		value.Write([]byte(string("\n" + username + " has joined the room\n# ")))
	}
	// handle user input
	handleCommands(c, username)
}

func handleCommands(c net.Conn, username string) {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		text := strings.TrimSpace(string(netData))
		if text == "" {
			c.Write([]byte(string("You need to say something...\n" + "# ")))
		} else if text == "QUIT" {
			c.Write([]byte(string("Thanks for chatting " + username + " !!\n")))
			c.Close()
			delete(connmap, username)
		} else {
			for _, value := range connmap {
				value.Write([]byte(string("\n" + username + ": " + text + "\n# ")))
			}
		}
	}
}

func handleConnection(c net.Conn) {
	// welcome user
	c.Write([]byte(string("Welcome to FlexCHAT !!\n")))
	// get user details
	username := createUser(c, "Please enter you username: ")
	// map username to connection
	connmap[username] = c
	// enter the chatroom
	enterRoom(c, username)
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		text := strings.TrimSpace(string(netData))
		if text == "QUIT" {
			c.Write([]byte(string("Thanks for chatting !!\n")))
			break
		}
	}
	c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
