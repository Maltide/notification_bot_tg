package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("before:", start.Before(time.Now()))
	fmt.Println("after:", start.After(time.Now()))
}
