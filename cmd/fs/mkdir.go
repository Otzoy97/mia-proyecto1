package fs

import "github.com/fatih/color"

//Mkdir ...
type Mkdir struct {
	id, path  string
	Row       int
	recursive bool
}

//AddOp ...
func (m *Mkdir) AddOp(s string, v interface{}) {
	if s == "id" {
		m.id = v.(string)
	} else if s == "path" {
		m.path = v.(string)
	} else if s == "p" {
		m.recursive = true
	}
}

//Validate ...
func (m *Mkdir) Validate() bool {
	flag := true
	if m.id == "" {
		flag = false
		color.New(color.FgHiYellow).Printf("Mkdirk: no se encontró id (%v)\n", m.Row)
	}
	if m.path == "" {
		flag = false
		color.New(color.FgHiYellow).Printf("Mkdirk: no se encontró path (%v)\n", m.Row)
	}
	if !flag {
		return false
	}
	return true
}

//Run crea un nuevo archivo
func (m *Mkdir) Run() {

}
