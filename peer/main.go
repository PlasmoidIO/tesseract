package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"share/common/packet"
	"share/peer/client"
	"strings"
)

func catch(err error) {
	if err != nil {
		log.Fatal(err)
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
		filesize := 815
		target := arr[1]
		fmt.Println("post result", cl.SendFile(filename, filesize, target))
	}
}
