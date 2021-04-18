package main

import (
	"bufio"
	"fmt"
	"os"
	"share/peer/client"
)

func inputLoop(c *client.CentralClient) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		err := c.RegisterUsername(text)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Printf("Username %s registered successfully\n", text)
	}
}

func main() {
	c := client.NewClient()
	go inputLoop(&c)
	c.Start()
}
