package main

import (
	"bufio"
	"fmt"
	"os"
	"share/common/packet"
	"share/peer/application"
	"strconv"
	"strings"
)

func main() {
	app := application.NewApplication()
	go app.Start()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " ")
		if len(arr) < 4 {
			continue
		}
		filename := arr[0]
		filesize, err := strconv.Atoi(arr[1])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		username := arr[2]
		addr := arr[3]

		p := packet.NewSendPacket(filename, filesize, username, addr)
		app.Client.WritePacket(&p)
	}
}
