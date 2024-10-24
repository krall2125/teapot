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
		for _, t := range(tokens) {
			fmt.Printf("[%s '%s' %d %d]\n", lexer.TT_ToStr(t.Typ), t.Lexeme, t.Line, t.Char);
		}
	}
}
