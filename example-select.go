package main

import "fmt"

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		var v int
		for {
			v = <-c
			fmt.Println(v)
			if v > 100 {
			  quit <- 1
			  break
			}
		}
		
	}()
	fibonacci(c, quit)
}
