package main

import (
	"bufio"
	"fmt"
	"log"
	"mia-proyecto1/analyzer"
	"os"
	"strings"
)

func main() {
	//line := []byte("mkdisk -path->/home/usr\n")
	in := bufio.NewReader(os.Stdin)
	auxLine := []byte{}
	aux := ""
	fmt.Println("Ready")
	lex := analyzer.Lexer{}
	par := analyzer.Parser{Lex: &lex}
	for {
	LeerEntrada:
		//Lee la entrada del usuario
		line, err := in.ReadBytes('\n')
		if err != nil {
			log.Fatalf("Error al leer entrada del usuario: %s", err)
		}
		//Convierte line a una cadena y
		//asi poder realizar operaciones entre cadenas
		aux = string(line)
		aux = strings.TrimSpace(aux)
		//Une la entrada del usuario a auxLine
		auxLine = append(auxLine, line...)
		//Si en line viene los caracteres '\*'
		//regresa a LeerEntrada para leer otra entrada del usuario
		if len(aux) > 2 && aux[len(aux)-2:] == "\\*" {
			goto LeerEntrada
		}
		//Analiza los bytes almacenados en auxLine
		lex.Line = auxLine
		lex.Scanner()
		par.Parser()
		for i := 0; i < len(par.Cmdlst); i++ {
			par.Cmdlst[i].Run()
		}
		//Reinicia auxLine
		auxLine = nil
	}
}
