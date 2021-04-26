package main

import (
	"bufio"
	"fmt"
	"os"
	"share/common/packet"
	"share/peer/client"
	"strings"
)

func catch(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func main() {
	cl := client.NewClient()
	go cl.Start()

	scanner := bufio.NewScanner(os.Stdin)
	cl.HandleSendRequest(func(p *packet.SendPacket) bool {
		fmt.Println("received this bitch")
		return true
	})
	fmt.Print("Username: ")
	if scanner.Scan() {
		catch(cl.RegisterUsername(scanner.Text()))
	}
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " ")
		if len(arr) < 2 {
			continue
		}
		filename := arr[0]
		target := arr[1]

		filesize, err := cl.SendFile(filename, target)
		if err != nil {
			fmt.Printf("Error sending file: %s\n", err)
		}
		if filesize >= 0 {
			fmt.Printf("File %s of size %d sent to %s successfully!\n", filename, filesize, target)
		} else {
			fmt.Printf("File %s sending to %s failed!\n", filename, target)
		}
	}
}
