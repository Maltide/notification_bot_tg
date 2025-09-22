package main

import (
	"fmt"
	"time"
)

func main() {

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "one"
		fmt.Println("one is written")
		c1 <- "one"
		fmt.Println("one is written")
	}()
	go func() {
		time.Sleep(1 * time.Second)
		c2 <- "two"
		fmt.Println("two is written")
	}()

	for range 3 {
		select {
		case msg := <-c1:
			time.Sleep(time.Millisecond)
			fmt.Println("received", msg)
		case msg := <-c2:
			time.Sleep(time.Millisecond)
			fmt.Println("received", msg)
		}
	}
}
