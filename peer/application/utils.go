package application

import "fmt"

func catch(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
