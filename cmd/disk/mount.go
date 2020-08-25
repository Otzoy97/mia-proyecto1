package cmdisk

import (
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"os"

	"github.com/fatih/color"
)

//Mount ...
type Mount struct {
	Row        int
	path, name string
	exec       byte
	Oplst      map[string]interface{}
}

//AddOp ...
func (m *Mount) AddOp(k string, v interface{}) {
	m.Oplst[k] = v
}

//Validate ...
func (m *Mount) Validate() bool {
	f := true
	m.exec = 'm'
	if len(m.Oplst) == 0 {
		//Se debe listar las particiones montadas
		m.exec = 'l'
		return true
	}
	if !cmd.ValidateOptions(&m.Oplst, "path") {
		color.New(color.FgHiYellow).Printf("Mount: path no se encontró (%v)\n", m.Row)
		f = false
	} else {
		m.path = m.Oplst["path"].(string)
	}
	if !cmd.ValidateOptions(&m.Oplst, "name") {
		color.New(color.FgHiYellow).Printf("Mount: name no se encontró (%v)\n", m.Row)
		f = false
	} else {
		m.name = m.Oplst["name"].(string)
	}
	if !f {
		color.New(color.FgHiRed, color.Bold).Println("Mount no se puede ejecutar")
		return false
	}
	return true
}

//Run almacena el path de un disco y el name de una partición
func (m *Mount) Run() {
	switch m.exec {
	case 'm':
		//Verifica que la partición exista
		//Verifica el disco
		if _, err := os.Stat(m.path); err != nil {
			color.New(color.FgHiYellow).Printf("Mount: el disco '%v' no existe o no se puede abrir (%v)\n", m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Mount fracasó")
			return
		}
		//Abre el disco y verifica los nombres de particiones
		file, err := os.Open(m.path)
		if err != nil {
			color.New(color.FgHiYellow).Printf("Mount: no se puede abrir el disco '%v' (%v)\n", m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Mount fracasó")
			return
		}
		//Recupera el mbr
		mbr := disk.Mbr{}
		mbr.ReadMbr(file)
		arr, _ := cmd.CreateArrPart(&mbr)
		//Si el arreglo es < 0, no hay nada que montar
		if len(arr) == 0 {
			color.New(color.FgHiYellow).Printf("Mount: no hay particiones en el disco '%v' (%v)\n", m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Mount fracasó")
			return
		}
		//Verifica que el nombre exista
		if !cmd.CheckNames(&arr, m.name) {
			color.New(color.FgHiYellow).Printf("Mount: '%v' no existe en el disco '%v' (%v)\n", m.name, m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Mount fracasó")
			return
		}
		//Almacena el path y el name
		if ok, id := disk.AddImg(m.path, m.name); ok {
			color.New(color.FgHiGreen, color.Bold).Printf("Mount: se montó '%v' con id '%v' (%v)\n", m.name, id, m.Row)
		} else {
			color.New(color.FgHiRed).Printf("Mount fracasó (%v)\n", m.Row)
		}
	case 'l':
		//Listar particiones
		disk.ListImg()
	}
}
