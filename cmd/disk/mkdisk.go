package cmdisk

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"os"
	"strings"
	"time"

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
func (m *Mkdisk) AddOp(key string, value interface{}) {
	m.Oplst[key] = value
}

//Validate ...
func (m *Mkdisk) Validate() bool {
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
		//Verifica que name tenga extensión dsk
		sCon := strings.Split(m.Oplst["name"].(string), ".")
		if len(sCon) == 2 && sCon[1] == "dsk" {
			m.name = m.Oplst["name"].(string)
		} else {
			color.New(color.FgHiYellow).Printf("Mkdisk: name debe tener extensión .dsk (%v)\n", m.Row)
			f = false
		}
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
			color.New(color.FgHiYellow).Printf("Mkdisk: unit debe ser 'k' o 'm' (%v)\n", m.Row)
			f = false
		}
	}
	if !f {
		color.New(color.FgHiRed, color.Bold).Println("Mkdisk no se puede ejecutar")
		return false
	}
	return true
}

//Run crea un disco
func (m *Mkdisk) Run() {
	//Crea el directorio
	if !m.createFolders() {
		color.New(color.FgHiRed, color.Bold).Println("Mkdisk fracasó")
		return
	}
	//Crea el mbr
	mbr := disk.Mbr{MbrTamanio: uint32(m.size * m.unit * 1024), MbrDiskSignature: uint32(rand.Intn(10000000))}
	//Guarda la fecha/hora
	tDec, _ := time.Now().GobEncode()
	copy(mbr.MbrFechaCreacion[:], tDec)
	//Crea el archivo
	file, err := os.Create(m.path + "/" + m.name)
	defer file.Close()
	if err != nil {
		color.New(color.FgHiYellow).Printf("Mkdisk: no se pudo crear el archivo '%v' (%v)\n%v\n", m.name, m.Row, err.Error())
		color.New(color.FgHiRed, color.Bold).Println("Mkdisk fracasó")
		return
	}
	bin := new(bytes.Buffer)
	arrB := make([]byte, 1024*m.unit)
	binary.Write(bin, binary.BigEndian, &arrB)
	for i := 0; i < m.size; i++ {
		file.Write(bin.Bytes())
	}
	if mbr.WriteMbr(file) {
		color.New(color.FgHiGreen, color.Bold).Printf("Mkdisk: se creó el disco '%v' (%v)\n", m.path+"/"+m.name, m.Row)
	} else {
		color.New(color.FgHiRed, color.Bold).Printf("Mkdisk fracasó (%v)\n", m.Row)
	}
}

//Verifica que cada carpeta del path exista. Si no existe la crea
func (m *Mkdisk) createFolders() bool {
	if err := os.MkdirAll(m.path, os.ModePerm); err != nil {
		color.New(color.FgHiYellow).Printf("Mkdisk: no se pudo crear el directorio '%v' (%v)\n%v\n", m.path, m.Row, err.Error())
		return false
	}
	return true
}
