package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		c := make(chan int)
		go Walk(t.Left, c)
		for {
			v, ok := <-c
			if !ok {
				break
			}
			ch <- v
		}
	}
	ch <- t.Value
	if t.Right != nil {
		c := make(chan int)
		go Walk(t.Right, c)
		for {
			v, ok := <-c
			if !ok {
				break
			}
			ch <- v
		}
	}
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if v1 != v2 || ok1 != ok2 {
			return false
		}
		if !ok1 {
			break
		}
	}
	return true
}

func main() {
	c := make(chan int)
	go Walk(tree.New(1), c)
	for {
		v, ok := <-c
		if !ok {
			break
		}
		fmt.Println(v)
	}
	if Same(tree.New(1), tree.New(1)) {
		fmt.Println("Same!")
	} else {
		fmt.Println("Houston, we have a problem!")
	}
	if Same(tree.New(1), tree.New(2)) {
		fmt.Println("Uh-oh!")
	} else {
		fmt.Println("Ok, they are different indeed!")
	}
}
