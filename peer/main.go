package main

import (
	"bufio"
	"fmt"
	"os"
	"share/common/packet"
	"share/peer/client"
	"strings"
)

var scanner *bufio.Scanner


func handleSendRequest(p *packet.SendPacket) bool {
	fmt.Printf("Send request received from %s of a file named: %s, with a size of %d bytes.\nAccept? [y,n] ", p.User, p.Filename, p.Size)
	if !scanner.Scan() {
		return false
	}
	for scanner.Scan() {
		input := scanner.Text()
		switch strings.ToLower(input) {
		case "y":
			return true
		case "n":
			return false
		}
	}

	return false
}

func registerUsername(c *client.CentralClient) {
	fmt.Print("Username: ")
	scanner = bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	if err := c.RegisterUsername(scanner.Text()); err == nil {
		fmt.Println("Registered successfully")
	} else {
		fmt.Printf("Error registering: %s\n", err)
	}
}

func main() {
	c := client.NewClient("/ip4/127.0.0.1/54061/filesharing/hfodFHDSOhsodfnONFoSHFP34u0")
	go registerUsername(&c)
	c.HandleSendRequest(handleSendRequest)
	c.Start()
}
