package main

import "fmt"

func main() {
	var a, b int
	for {
		_, e := fmt.Scanln(&a, &b)
		if e != nil {
			break
		}
		fmt.Println(a + b)
	}
}
