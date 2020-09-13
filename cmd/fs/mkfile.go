package fs

import (
	"mia-proyecto1/disk"
	"mia-proyecto1/lwh"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

//Mkfile ...
type Mkfile struct {
	Row            int
	id, path, cont string
	recursive      bool
	size           int32
}

//AddOp ...
func (m *Mkfile) AddOp(s string, v interface{}) {
	if s == "id" {
		m.id = v.(string)
	} else if s == "path" {
		m.path = v.(string)
	} else if s == "p" {
		m.recursive = true
	} else if s == "size" {
		conv, _ := strconv.Atoi(v.(string))
		m.size = int32(conv)
	} else if s == "cont" {
		m.cont = v.(string)
	}
}

//Validate ...
func (m *Mkfile) Validate() bool {
	flag := true
	if !flag {
		return false
	}
	return true
}

//Run crea un nuevo archivo
func (m *Mkfile) Run() {
	//Verifica que exista un usuario logeado
	if !lwh.IsActive() {
		color.New(color.FgHiYellow).Printf("Mkfile: no hay sesión activa (%v)\n", m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
		return
	}
	//Verifica que la partición esté montada
	path, name := disk.FindImg(m.id)
	if path == "" && name == "" {
		color.New(color.FgHiYellow).Printf("Mkfile: partición con id '%v' no está montada (%v)\n", m.id, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
		return
	}
	//Monta el disco
	lwh.MountVDisk(path, name)
	//Verifica que el archivo no exista
	var avd lwh.Avd
	//Recupera el directorio raiz
	avd.ReadAvd(0)
	//Busca coincidencia con el path
	_, tip := avd.Find(m.path)
	switch tip {
	case 0, 1:
		//Ya existe el directorio/archivo
		color.New(color.FgHiYellow).Printf("Mkfile: '%v' ya existe (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
	case 2:
		//No existe el archivo
		//Verifica si es necesario utilizar recursión
		if m.validateRecursion(&avd) && !m.recursive {
			color.New(color.FgHiYellow).Printf("Mkfile: '%v' no se puede crear, especifique recursión (%v)\n", m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
		} else {
			//Se puede crear los archivos
			path := strings.Split(path, "/")
			path = path[1:]
			flag, res := avd.CreateDir("/", path)
			if !flag {
				switch res {
				case 0:
					color.New(color.FgHiYellow).Printf("Mkfile: '%v' no se creó, ya existe (%v)\n", m.path, m.Row)
					color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
				case 2:
					color.New(color.FgHiYellow).Printf("Mkfile: '%v' no se puede crear (%v)\n", m.path, m.Row)
					color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
				case 3:
					color.New(color.FgHiYellow).Printf("Mkfile: '%v' no es una dirección válida (%v)\n", m.path, m.Row)
					color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
				case 4:
					color.New(color.FgHiYellow).Printf("Mkfile: '%v' no se puede leer el disco (%v)\n", m.path, m.Row)
					color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
				}
			} else {
				color.New(color.FgHiGreen, color.Bold).Printf("Mkfile: '%v' creado exitosamente\n\n", m.path)
			}
		}
	case 3:
		//No es directorio
		color.New(color.FgHiYellow).Printf("Mkfile: '%v' no es un directorio válido (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
	case 4:
		//No pudo leerse alguna estructura
		color.New(color.FgHiYellow).Printf("Mkfile: hubo un error al leer el disco (%v)\n", m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkfile fracasó", "")
	}
	//Desmonta el disco
	lwh.UnmountVDisk()
}

//validateRecursion verifica si es necesaria la recursión
func (m *Mkfile) validateRecursion(root *lwh.Avd) bool {
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
