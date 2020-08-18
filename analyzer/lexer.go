package analyzer

import (
	"fmt"
	"unicode/utf8"
)

const eof = 0

var tokQueue []*Token

//Token ...
type Token struct {
	lex     interface{}
	row     int
	col     int
	tokname string
}

//Lexer ...
type Lexer struct {
	Row, Col int
	Line     []byte
	Peek     rune
}

var tokNames = []string{"mkdisk",
	"rmdisk",
	"fdisk",
	"mount",
	"unmount",
	"exec",
	"rep",
	"-path",
	"-size",
	"-fit",
	"-unit",
	"-type",
	"-delete",
	"-name",
	"-add",
	"-id",
	"bf",
	"ff",
	"wf",
	"m",
	"k",
	"b",
	"p",
	"e",
	"l",
	"fast",
	"full",
	"mbr",
	"disk"}

func (x *Lexer) next() rune {
	if x.Peek != eof {
		r := x.Peek
		x.Peek = eof
		return r
	}
	if len(x.Line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.Line)
	x.Line = x.Line[size:]
	x.Col++
	if c == utf8.RuneError && size == 1 {
		return x.next()
	}
	return c
}

//Scanner ...
func (x *Lexer) Scanner() {
	state := 0
	stringRec := ""
	numberRec := 0
	for {
		c := x.next()
		if c == eof {
			return
		}
		switch state {
		case 0:
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				stringRec = string(c)
				state = 1
			} else if c >= '0' && c <= '9' {
				numberRec = int(c) - '0'
				state = 3
			} else if c == '-' {
				stringRec = string(c)
				state = 4
			} else if c == '\\' {
				state = 8
			} else if c == '#' {
				state = 10
			} else if c == '"' {
				stringRec = ""
				state = 2
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				x.Row++
				x.Col = 0
				continue
			} else {
				fmt.Printf("unrecognized character %q (%v, %v)", c, x.Row, x.Col)
			}
		case 1:
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
				stringRec += string(c)
			} else {
				if c != eof {
					x.Peek = c
				}
				x.reservada(stringRec)
			}
		case 2:
			if c == '"' {
				if c != eof {
					x.Peek = c
				}
				tokQueue = append(tokQueue, &Token{lex: stringRec, row: x.Row, col: x.Col, tokname: "cadena"})
			} else if c != '\n' {
				stringRec += string(c)
			} else {
				fmt.Printf("unclosed string (%v, %v)", x.Row, x.Col)
			}
		case 3:
			if c >= '0' && c <= '9' {
				numberRec *= 10
				numberRec += int(c) - '0'
			} else {
				if c != eof {
					x.Peek = c
				}
				tokQueue = append(tokQueue, &Token{lex: numberRec, row: x.Row, col: x.Col, tokname: "numero"})
			}
		case 4:
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				stringRec += string(c)
				state = 5
			} else if c >= '0' && c <= '9' {
				numberRec = (int(c) - '0') * -1
				state = 7
			} else if c == '>' {
				tokQueue = append(tokQueue, &Token{lex: "->", row: x.Row, col: x.Col, tokname: "asignacion"})
			} else {
				fmt.Printf("unrecognized character %q (%v, %v)", c, x.Row, x.Col)
			}
		case 5:
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				stringRec += string(c)
			} else {
				if c != eof {
					x.Peek = c
				}
				x.reservada(stringRec)
			}
		case 7:
			if c >= '0' && c <= '9' {
				numberRec *= 10
				numberRec = (int(c) - '0') * -1
			} else {
				if c != eof {
					x.Peek = c
				}
				tokQueue = append(tokQueue, &Token{lex: numberRec, row: x.Row, col: x.Col, tokname: "numero"})
			}
		case 8:
			if c == '*' {
				state = 9
			} else {
				if c != eof {
					x.Peek = c
				}
				state = 0
			}
		case 9, 10:
			if c == '\n' {
				stringRec = ""
				numberRec = 0
				state = 0
				x.Row++
				x.Col = 0
			}
		default:
			return
		}
	}
}

func (x *Lexer) reservada(s string) {
	for _, v := range tokNames {
		if v == s {
			tokQueue = append(tokQueue, &Token{lex: s, row: x.Row, col: x.Col, tokname: v})
		}
	}
}
