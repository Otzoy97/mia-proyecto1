package fs

import (
	"mia-proyecto1/disk"
	"mia-proyecto1/lwh"
	"strings"

	"github.com/fatih/color"
)

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
	//Verifica que exista un usuario logeado
	if !lwh.IsActive() {
		color.New(color.FgHiYellow).Printf("Mkdir: no hay sesión activa (%v)\n", m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkdir fracasó", "")
		return
	}
	//Verifica que la partición esté montada
	path, name := disk.FindImg(m.id)
	if path == "" && name == "" {
		color.New(color.FgHiYellow).Printf("Mkdir: partición con id '%v' no está montada (%v)\n", m.id, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkdir fracasó", "")
		return
	}
	//Monta el disco
	lwh.MountVDisk(path, name)
	//Verifica que el directorio no exista
	var avd lwh.Avd
	//Recupera el directorio raiz
	avd.ReadAvd(0)
	//Busca coincidencia con el path
	_, tip := avd.Find(m.path)
	switch tip {
	case 0, 1:
		//Ya existe el directorio/archivo
		color.New(color.FgHiYellow).Printf("Mkdir: '%v' ya existe (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkdir fracasó", "")
	case 2:
		//No existe el archivo
		//Verifica si es necesario utilizar recursión
		if m.validateRecursion(&avd) && !m.recursive {
			color.New(color.FgHiYellow).Printf("Mkdir: '%v' no se puede crear, especifique recursión (%v)\n", m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Mkdir fracasó", "")
		} else {
			//Se puede crear los archivos
			path := strings.Split(m.path, "/")
			path = path[1:]
			root := "/"
			avd.CreateDir(&root, path)
			//Sobreescribe el avd
			avd.WriteAvd(0)
			color.New(color.FgHiBlue, color.Bold).Println("Mkdir ha finalizado", "")
		}
	case 3:
		//No es directorio
		color.New(color.FgHiYellow).Printf("Mkdir: '%v' no es un directorio válido (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkdir fracasó", "")
	case 4:
		//No pudo leerse alguna estructura
		color.New(color.FgHiYellow).Printf("Mkdir: hubo un error al leer el disco (%v)\n", m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkdir fracasó", "")
	}
	//Desmonta el disco
	lwh.UnmountVDisk()
}

//validateRecursion verifica si es necesaria la recursión
func (m *Mkdir) validateRecursion(root *lwh.Avd) bool {
	//Se asume que el path es válido
	pathSplited := strings.Split(m.path, "/")
	//Retira el último directorio
	pathSplited = pathSplited[1 : len(pathSplited)-1]
	//Vuelve a unir todo el path
	auxPath := "/"
	for _, s := range pathSplited {
		auxPath += s
	}
	//Busca el directorio
	_, tipe := root.Find(auxPath)
	//tipe debe ser un directorio
	if tipe == 0 {
		return false
	}
	return true
}
