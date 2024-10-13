package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int) // an unbuffered channel
	go func(ch chan<- int, x int) {
		time.Sleep(time.Second)
		
		ch <- x*x // 9 is sent
	}(c, 3)
	done := make(chan struct{})
	go func(ch <-chan int) {
		
		n := <-ch
		fmt.Println(n) // 9
		
		time.Sleep(time.Second)
		done <- struct{}{}
	}(c)
	
	<-done
	fmt.Println("bye")
}
