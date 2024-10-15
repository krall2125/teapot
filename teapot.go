package main

import (
	"fmt"
	"os"
	"teapot/lexer"
)

func main() {
	args := os.Args[1:]

	for _, v := range(args) {
		tokens := lexer.LexFile(v);
		fmt.Printf("Tokens: %v\n", tokens);
	}
}
