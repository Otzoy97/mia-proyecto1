package rep

import (
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

//Rep ...
type Rep struct {
	Oplst                  map[string]interface{}
	path, id, ruta, nombre string
	exec                   bool
	Row                    int
}

//AddOp ...
func (m *Rep) AddOp(k string, v interface{}) {
	m.Oplst[k] = v
}

//Validate ...
func (m *Rep) Validate() bool {
	f := true
	m.exec = false
	bRuta := cmd.ValidateOptions(&m.Oplst, "ruta")
	if cmd.ValidateOptions(&m.Oplst, "nombre") {
		//Si el nombre es igual a ls o file, bRuta sería obligatorio
		if m.Oplst["nombre"].(string) == "ls" || m.Oplst["nombre"].(string) == "file" {
			if bRuta {
				m.ruta = m.Oplst["ruta"].(string)
				m.exec = true
			} else {
				color.New(color.FgHiYellow).Printf("Rep: ruta no se encontró (%v)\n", m.Row)
				f = false
			}
		}
	} else {
		color.New(color.FgHiYellow).Printf("Rep: nombre no se encontró (%v)\n", m.Row)
		f = false
	}
	if !cmd.ValidateOptions(&m.Oplst, "path") {
		color.New(color.FgHiYellow).Printf("Rep: path no se encontró (%v)\n", m.Row)
		f = false
	}
	if !cmd.ValidateOptions(&m.Oplst, "id") {
		color.New(color.FgHiYellow).Printf("Rep: id no se encontró (%v)\n", m.Row)
		f = false
	}
	if !f {
		color.New(color.FgHiRed, color.Bold).Println("Rep no se puede ejecutar")
		return false
	}
	m.nombre = m.Oplst["nombre"].(string)
	m.path = m.Oplst["path"].(string)
	m.id = m.Oplst["id"].(string)
	return true
}

//Run ...
func (m *Rep) Run() {
	//Recupera el mbr
	mbr := disk.Mbr{}
	//Recupera el path del disco
	path, name := disk.FindImg(m.id)
	if path == "" {
		//No existe la partición
		color.New(color.FgHiYellow).Printf("Rep: '%v' no está montada (%v)\n", m.id, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rep fracasó")
		return
	}
	//Intenta abrir el disco
	diskFile, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	defer diskFile.Close()
	//Verifica que no haya error
	if err != nil {
		color.New(color.FgHiYellow).Printf("Rep: el disco '%v' no existe o no se puede abrir (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rep fracasó")
		return
	}
	//Recupera el mbr del disco
	if !mbr.ReadMbr(diskFile) {
		color.New(color.FgHiYellow).Printf("Rep: no se pudo recuperar el mbr del disco '%v' (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rep fracasó")
		return
	}
	//Crea el directorio de path
	if err := os.MkdirAll(m.path, os.ModePerm); err != nil {
		color.New(color.FgHiYellow).Printf("Rep: no se pudo crear el directorio '%v' (%v)\n(%v)\n", m.path, m.Row, err.Error())
		color.New(color.FgHiRed).Println("Rep fracasó")
		return
	}
	//Crea el archivo qu contendra el texto para generar el grafo
	gFile, _ := os.Create("doFile.dot")
	switch m.nombre {
	case "mbr":
		//Construye un texto para el reporte del MBR
		b := m.CreateMBR(&mbr)
		//Escribe el archivo dot
		gFile.Write(b)
		gFile.Close()
	case "disk":
		//Construye el text para el reporte del DISK
		b := m.CreateDisk(&mbr)
		//Escribe el archivo dot
		gFile.Write(b)
		gFile.Close()
	}
	//Crea el nombre para el reporte
	fileName := m.nombre + "_" + strings.Split(filepath.Base(path), ".")[0] + "_" + name
	//Crea el grafo
	graphName := m.path + "/" + fileName + ".png"
	cmd := exec.Command("dot", "-Tpng", "doFile.dot", "-o", graphName)
	if err := cmd.Run(); err != nil {
		color.New(color.FgHiYellow).Printf("Rep: graphviz no pudo generar la imagen (%v)\n%v", m.Row, err.Error())
		color.New(color.FgHiRed, color.Bold).Println("Rep fracasó")
	} else {
		os.Remove("doFile.dot")
		color.New(color.FgHiGreen, color.Bold).Printf("Rep generó '%v' en '%v' (%v)\n", fileName, m.path, m.Row)
	}
}
