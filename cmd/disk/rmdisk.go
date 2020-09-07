package cmdisk

import (
	"bufio"
	"mia-proyecto1/cmd"
	"os"
	"strings"

	"github.com/fatih/color"
)

//Rmdisk ...
type Rmdisk struct {
	size, Row int
	path      string
	Oplst     map[string]interface{}
}

//AddOp ...
func (m *Rmdisk) AddOp(k string, v interface{}) {
	m.Oplst[k] = v
}

//Validate ...
func (m *Rmdisk) Validate() bool {
	if !cmd.ValidateOptions(&m.Oplst, "path") {
		color.New(color.FgHiYellow).Printf("Rmdisk: path no se encontró (%v)\n", m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rmdisk no se puede ejecutar")
		return false
	}
	m.path = m.Oplst["path"].(string)
	return true
}

//Run elimina un disco
//TODO: Validar que no se elimine un disco que tenga montada una partición
func (m *Rmdisk) Run() {
	//Verifica si el disco existe.
	if _, err := os.Stat(m.path); err != nil {
		color.New(color.FgHiYellow).Printf("Rmdisk: el disco no existe %v\n", m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rmdisk fracasó")
		return
	}
	//Solicita confirmación para eliminar el disco
	in := bufio.NewReader(os.Stdin)
ReEntry:
	color.New(color.FgHiBlue).Printf("¿Desea eliminar el disco '%v'? [s/n] ", m.path)
	txt, err := in.ReadString('\n')
	if err != nil {
		color.New(color.FgHiRed, color.Bold).Println("Error al leer la entrada del usuario.")
		color.New(color.FgHiRed, color.Bold).Println("Rmdisk fracasó")
		return
	}
	//Se asegura que la entrada sea /s/ o /n/
	txt = strings.ToLower(strings.TrimSpace(txt))
	if txt != "s" && txt != "n" {
		goto ReEntry
	}
	if txt == "n" {
		color.New(color.FgHiBlue, color.Bold).Printf("Rmdisk: el disco '%v' no se eliminó\n", m.path)
		return
	}
	//Elimina el disco
	if err := os.Remove(m.path); err != nil {
		color.New(color.FgHiYellow).Printf("Rmdisk: no se pudo eliminar el disco '%v'\n%v", m.path, err.Error())
		return
	}
	color.New(color.FgHiGreen, color.Bold).Printf("Rmdisk: se eliminó el disco '%v'\n", m.path)
}
