package main

import (
	"fmt"
	"time"
)

func main() {
	announce("hello ", 5)
}

func announce(message string, delay int) {
	go func() {
		fmt.Println(" message is %v ", message)
		time.Sleep(delay * time.Second)
	}()
}
