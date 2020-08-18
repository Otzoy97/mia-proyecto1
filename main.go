package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"miap1/analyzer"
	"os"
)

func main() {
	fmt.Println("ada")
	const src = `mkdisk -path->/home/usr/ -name->disco1 -add->-80`
	in := bufio.NewReader(os.Stdin)
	for {
		if _, err := os.Stdout.WriteString(""); err != nil {
			log.Fatalf("WriteString: %s", err)
		}
		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}
		lex := analyzer.Lexer{Line: line}
		lex.Scanner()
	}

}
