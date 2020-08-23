package analyzer

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

const eof = 0

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
	tokQueue []*Token
}

var tokNames map[string]string = map[string]string{"mkdisk": "mkdisk",
	//"rmdisk":          "rmdisk",
	"fdisk":           "fdisk",
	"mount":           "mount",
	"unmount":         "unmount",
	"exec":            "exec",
	"rep":             "rep",
	"mkfs":            "mkfs",
	"login":           "login",
	"logout":          "logout",
	"mkgrp":           "mkgrp",
	"mkusr":           "mkusr",
	"mkfile":          "mkfile",
	"mkdir":           "mkdir",
	"loss":            "loss",
	"recovery":        "recovery",
	"pause":           "pause",
	"-nombre":         "nombre",
	"-path":           "path",
	"-size":           "size",
	"-fit":            "fit",
	"-unit":           "unit",
	"-type":           "type",
	"-delete":         "delete",
	"-name":           "name",
	"-add":            "add",
	"-id":             "id",
	"-ruta":           "ruta",
	"-p":              "p",
	"-cont":           "cont",
	"-usr":            "usr",
	"-pwd":            "pwd",
	"-grp":            "grp",
	"-tipo":           "tipo",
	"bf":              "cadena", //"opFit",
	"ff":              "cadena", //"opFit",
	"wf":              "cadena", //"opFit",
	"m":               "cadena", //"opUnit",
	"k":               "cadena", //"opUnit",
	"b":               "cadena", //"opUnit",
	"p":               "cadena", //"opType",
	"e":               "cadena", //"opType",
	"l":               "cadena", //"opType",
	"fast":            "cadena", //"opDel",
	"full":            "cadena", //"opDel",
	"mbr":             "cadena", //"opRep",
	"disk":            "cadena", //"opRep",
	"sb":              "cadena", //"opRep",
	"bm_arbdir":       "cadena", //"opRep",
	"bm_detdir":       "cadena", //"opRep",
	"bm_inode":        "cadena", //"opRep",
	"bm_block":        "cadena", //"opRep",
	"bitacora":        "cadena", //"opRep",
	"directorio":      "cadena", //"opRep",
	"tree_file":       "cadena", //"opRep",
	"tree_directorio": "cadena", //"opRep",
	"tree_complete":   "cadena", //"opRep",
	"ls":              "cadena"} //"opRep"

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
	x.tokQueue = nil
	state := 0
	stringRec := ""
	numberRec := 0
	c := x.next()
	for c != eof {
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
			} else if c == '/' {
				stringRec = string(c)
				state = 11
			} else if c == '#' {
				state = 10
			} else if c == '"' {
				stringRec = ""
				state = 2
			} else if c == ' ' || c == '\t' || c == '\r' {
			} else if c == '\n' {
				x.Row++
				x.Col = 0
			} else {
				color.New(color.FgHiYellow).Printf("Caracter no reconocido %q (%v, %v)\n", c, x.Row, x.Col)
			}
		case 1:
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '.' {
				stringRec += string(c)
			} else {
				if c != eof {
					x.Peek = c
				}
				x.reservada(stringRec)
				state = 0
			}
		case 2:
			if c == '"' {
				state = 0
				x.tokQueue = append(x.tokQueue, &Token{lex: stringRec, row: x.Row, col: x.Col - len(stringRec), tokname: "cadena"})
			} else if c != '\n' {
				stringRec += string(c)
			} else {
				state = 0
				if c != eof {
					x.Peek = c
				}
				color.New(color.FgHiYellow).Printf("Cadena sin cerrar (%v, %v)\n", x.Row, x.Col)
			}
		case 3:
			if c >= '0' && c <= '9' {
				numberRec *= 10
				numberRec += int(c) - '0'
			} else {
				if c != eof {
					x.Peek = c
				}
				state = 0
				x.tokQueue = append(x.tokQueue, &Token{lex: numberRec, row: x.Row, col: x.Col - len(strconv.Itoa(numberRec)), tokname: "numero"})
			}
		case 4:
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				stringRec += string(c)
				state = 5
			} else if c >= '0' && c <= '9' {
				numberRec = (int(c) - '0') * -1
				state = 7
			} else if c == '>' {
				state = 0
				x.tokQueue = append(x.tokQueue, &Token{lex: "->", row: x.Row, col: x.Col - len(stringRec), tokname: "asignacion"})
			} else {
				state = 0
				if c != eof {
					x.Peek = c
				}
				color.New(color.FgHiYellow).Printf("Caracter no reconocido %q (%v, %v)\n", c, x.Row, x.Col)
			}
		case 5:
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				stringRec += string(c)
			} else {
				if c != eof {
					x.Peek = c
				}
				state = 0
				x.reservada(stringRec)
			}
		case 7:
			if c >= '0' && c <= '9' {
				numberRec *= 10
				numberRec += (int(c) - '0') * -1
			} else {
				if c != eof {
					x.Peek = c
				}
				state = 0
				x.tokQueue = append(x.tokQueue, &Token{lex: numberRec, row: x.Row, col: x.Col - len(strconv.Itoa(numberRec)), tokname: "numero"})
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
		case 11:
			if (c < 'a' && c > 'z') || (c < 'A' && c > 'Z') || (c < '0' && c > '9') || c == '-' || c == '_' || c == 'ñ' || c == 'Ñ' {
				color.New(color.FgHiYellow).Printf("Caracter no reconocido %q (%v, %v)\n", c, x.Row, x.Col)
				if c != eof {
					x.Peek = c
				}
				state = 0
			} else {
				stringRec += string(c)
				state = 12
			}
		case 12:
			if c == '/' {
				state = 11
				stringRec += string(c)
			} else if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '-' && c != '_' && c != 'ñ' && c != 'Ñ' && c != '.' {
				if c != eof {
					x.Peek = c
				}
				x.tokQueue = append(x.tokQueue, &Token{lex: stringRec, row: x.Row, col: x.Col - len(stringRec), tokname: "cadena"})
				stringRec = ""
				state = 0
			} else {
				stringRec += string(c)
				state = 12
			}
		}
		c = x.next()
	}
	//Agrega un $ al final del stream de tokens
	x.tokQueue = append(x.tokQueue, &Token{lex: "$", row: x.Row, col: x.Col + 1, tokname: "$"})
}

func (x *Lexer) reservada(s string) {
	for k, v := range tokNames {
		if k == strings.ToLower(s) {
			x.tokQueue = append(x.tokQueue, &Token{lex: strings.ToLower(s), row: x.Row, col: x.Col - len(s), tokname: v})
			return
		}
	}
	x.tokQueue = append(x.tokQueue, &Token{lex: s, row: x.Row, col: x.Col - len(s), tokname: "cadena"})
}
