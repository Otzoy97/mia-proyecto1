package main

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

func main() {
	const src = `mkdisk -path->/home/usr/ -name->disco1 -add->-80 \*`

	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "example"
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '-' && i == 0 || unicode.IsLetter(ch) || ch == '/' || unicode.IsDigit(ch) && i > 0
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}

}

func lex() {
	for {

	}
}
