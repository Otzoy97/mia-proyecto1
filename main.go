package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mia-proyecto1/analyzer"
	"os"
)

func main() {
	//src := []byte("mkdisk -path->\"/home/usr/\" -name->disco1 -add->-80\n")
	in := bufio.NewReader(os.Stdin)
	fmt.Println("Preaparado")
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
	//lex := analyzer.Lexer{Line: src}
	//lex.Scanner()

}
