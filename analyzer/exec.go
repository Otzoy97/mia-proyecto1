package analyzer

import (
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"golang.org/x/tools/godoc/util"
)

//Exec ...
type Exec struct {
	Row  int
	path string
}

//AddOp ...
func (m *Exec) AddOp(s string, v interface{}) {
	if s == "path" {
		m.path = v.(string)
	}
}

//Validate ...
func (m *Exec) Validate() bool {
	if m.path == "" {
		return false
	}
	return true
}

//Run ...
func (m *Exec) Run() {
	//Verificar si el archivo exite
	if _, err := os.Stat(m.path); err != nil {
		color.New(color.FgHiYellow).Printf("Exec: el archivo '%v' no existe o no se puede abrir (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Exec fracasó")
		return
	}
	//No es un archivo de texto
	// if !util.IsTextFile(f, m.path) {
	// 	color.New(color.FgHiYellow).Printf("Exec: el archivo '%v' no es un archivo de texto (%v)\n", m.path, m.Row)
	// 	color.New(color.FgHiRed, color.Bold).Println("Exec fracasó")
	// 	return
	// }
	//Abre el archivo
	content, err := ioutil.ReadFile(m.path)
	//Verifica que es un archivo de texto
	if !util.IsText(content) {
		color.New(color.FgHiYellow).Printf("Exec: el archivo '%v' no es un archivo de texto (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Exec fracasó")
		return
	}
	if err != nil {
		color.New(color.FgHiYellow).Printf("Exec: ocurrió un error al leer el contenido del archivo '%v' (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Exec fracasó")
		return
	}
	//Recupera el contenido del archivo
	lex := Lexer{Line: content}
	par := Parser{Lex: &lex}
	//Ejecuta el análisis
	lex.Scanner()
	par.Parser()
	for i := 0; i < len(par.Cmdlst); i++ {
		if par.Cmdlst[i].Validate() {
			par.Cmdlst[i].Run()
		}
	}
}
