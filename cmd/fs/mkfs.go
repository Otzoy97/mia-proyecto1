package mkfs

import (
	"mia-proyecto1/disk"
	"mia-proyecto1/lwh"
	"os"
	"strings"

	"github.com/fatih/color"
)

//Mkfs ...
type Mkfs struct {
	id   string
	tipo byte
	Row  int
}

//AddOp ...
func (m *Mkfs) AddOp(s string, v interface{}) {
	if s == "id" {
		m.id = v.(string)
	} else if s == "type" {
		m.tipo = v.(byte)
	}
}

//Validate ....
func (m *Mkfs) Validate() bool {
	if m.tipo == 0 {
		m.tipo = 'u'
	}
	if m.id == "" {
		return false
	}
	return true
}

//Run ...
func (m *Mkfs) Run() {
	// Busca que la partición esté montada
	// Recupera el archivo del disco y la partición especificada
	// Solo trabaja sobre particiones primarias
	path, name := disk.FindImg(m.id)
	if path == "" {
		//No existe la partición montada
		color.New(color.FgHiRed, color.Bold).Printf("Mkfs fracasó (%v)\n", m.Row)
		return
	}
	//Intenta abrir el disco
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: el disco '%v' no se puede abrir o no existe (%v)\n", path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkfs fracasó")
		return
	}
	m.setStructs(file, name, path)
}

//Crea y configura las estructuras para el sistema de archivoss
func (m *Mkfs) setStructs(file *os.File, name string, path string) {
	//Recupera la información del mbr y la partición
	var mbr disk.Mbr
	mbr.ReadMbr(file)
	parArr, _ := mbr.CreateArrPart()
	//Se asume que la partición DEBE existir
	//Recupera la partición
	par := parArr.Find(name)
	//Realiza un split para recuperar el nombre virtual del disco
	strSplt := strings.Split(path, "/")
	hdName := strings.Split(strSplt[len(strSplt)-1], ".")[0]
	//Configura el superboot
	var sb lwh.Superboot
	sb.New(par, hdName)
	//Crea el AVD raíz
	//Crea un detalle de directorio para la raiz
	//Crea el archivo user.txt para el directorio raiz
	//TODO: Realizar las funciones para determinar los bitmaps disponibles
}
