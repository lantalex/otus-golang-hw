package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	input := "Hello, OTUS!"
	reversed := reverse.String(input)
	fmt.Println(reversed)
}
