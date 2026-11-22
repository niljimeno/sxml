package main

import (
	"fmt"
	"os"
	"sxml/tokenizer"
)

func main() {
	template, _ := os.ReadFile("./test.scm")
	fmt.Printf("%s\n", template)

	tokens, _ := tokenizer.Tokenize(template)
	fmt.Printf("%s\n", tokens)
}
