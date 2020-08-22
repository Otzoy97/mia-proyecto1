package cmdisk

import (
	"fmt"
	"mia-proyecto1/cmd"

	"github.com/fatih/color"
)

//Mkdisk ...
type Mkdisk struct {
	size, Row int
	path      string
	name      string
	unit      int
	Oplst     map[string]interface{}
}

//AddOp ...
func (m Mkdisk) AddOp(key string, value interface{}) {
	m.Oplst[key] = value
}

//Validate ...
func (m Mkdisk) Validate() bool {
	f := true
	m.unit = 1024
	if !cmd.ValidateOptions(&m.Oplst, "path") {
		color.New(color.FgHiYellow).Printf("Mkdisk: path no se encontró (%v)\n", m.Row)
		f = false
	} else {
		m.path = m.Oplst["path"].(string)
	}
	if !cmd.ValidateOptions(&m.Oplst, "name") {
		color.New(color.FgHiYellow).Printf("Mkdisk: name no se encontró (%v)\n", m.Row)
		f = false
	} else {
		m.name = m.Oplst["name"].(string)
	}
	if !cmd.ValidateOptions(&m.Oplst, "size") {
		color.New(color.FgHiYellow).Printf("Mkdisk: size no se encontró (%v)\n", m.Row)
		f = false
	} else {
		m.size = m.Oplst["size"].(int)
		if m.size <= 0 {
			color.New(color.FgHiYellow).Printf("Mkdisk: size debe ser mayor a cero (%v)\n", m.Row)
			f = false
		}
	}
	if cmd.ValidateOptions(&m.Oplst, "unit") {
		switch m.Oplst["unit"].(string) {
		case "k":
			m.unit = 1
		case "m":
			m.unit = 1024
		default:
			color.New(color.FgHiYellow).Println("Mkdisk: unit deb ser 'k' o 'm'")
			f = false
		}
	}
	if !f {
		color.New(color.FgHiRed, color.Bold).Println("Mkdisk no se puede ejecutar")
		return false
	}
	return true
}

//Run ...
func (m Mkdisk) Run() {
	for k, v := range m.Oplst {
		fmt.Printf("%v -> %v\n", k, v)
	}
}
