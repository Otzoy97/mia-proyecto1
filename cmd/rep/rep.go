package rep

import (
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"mia-proyecto1/lwh"
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
	diskPath, namePath     string
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
	if cmd.ValidateOptions(&m.Oplst, "name") {
		//Si el nombre es igual a ls o file, bRuta sería obligatorio
		if m.Oplst["name"].(string) == "ls" || m.Oplst["name"].(string) == "file" {
			if bRuta {
				m.ruta = m.Oplst["ruta"].(string)
				m.exec = true
			} else {
				color.New(color.FgHiYellow).Printf("Rep: ruta no se encontró (%v)\n", m.Row)
				f = false
			}
		}
	} else {
		color.New(color.FgHiYellow).Printf("Rep: name no se encontró (%v)\n", m.Row)
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
	m.nombre = m.Oplst["name"].(string)
	m.path = m.Oplst["path"].(string)
	m.id = m.Oplst["id"].(string)
	return true
}

//Run ...
func (m *Rep) Run() {
	//Recupera el mbr
	mbr := disk.Mbr{}
	//Recupera el path del disco
	m.diskPath, m.namePath = disk.FindImg(m.id)
	if m.diskPath == "" {
		//No existe la partición
		color.New(color.FgHiYellow).Printf("Rep: '%v' no está montada (%v)\n", m.id, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rep fracasó")
		return
	}
	//Intenta abrir el disco
	diskFile, err := os.OpenFile(m.diskPath, os.O_RDWR, os.ModePerm)
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
	//Crea el nombre para el reporte
	fileName := m.nombre + "_" + strings.Split(filepath.Base(m.diskPath), ".")[0] + "_" + m.namePath
	//Recupera el contenido del archivo que generará el reporte
	gContent := m.createCont(&mbr)
	//Decide qué tipo de reporte generará (imagen, archivo de texto)
	if m.nombre == "bm_arbdir" || m.nombre == "bm_detdir" || m.nombre == "bm_inode" || m.nombre == "bm_block" || m.nombre == "bitacora" {
		newFile, _ := os.Create(m.path + "/" + fileName + ".txt")
		newFile.Write(gContent)
		newFile.Close()
		color.New(color.FgHiGreen, color.Bold).Printf("Rep generó '%v' en '%v' (%v)\n", fileName, m.path, m.Row)
		return
	}
	//Crea el archivo qu contendra el texto para generar el grafo
	gFile, _ := os.Create("doFile.dot")
	//Escribe el gContent en gFile
	gFile.Write(gContent)
	gFile.Close()
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

//createCont crea el contenido del archivo de texto para generar el reporte
func (m *Rep) createCont(mbr *disk.Mbr) []byte {
	switch m.nombre {
	case "mbr":
		//Construye un texto para el reporte del MBR
		return m.CreateMBR(mbr)
	case "disk":
		//Construye el text para el reporte del DISK
		return m.CreateDisk(mbr)
	case "sb":
		//Monta el sistema de archivos
		lwh.MountVDisk(m.diskPath, m.namePath)
		//Construye el texto para el reporte del SUPERBOOT
		b := m.CreateSB(lwh.GetSuperboot())
		//Desmonta el sistema de archivos
		lwh.UnmountVDisk()
		return b
	case "bm_arbdir":
		//Monta el sistema de archivos
		lwh.MountVDisk(m.diskPath, m.namePath)
		//Recupera el bitmap del arbol de directorio
		b := lwh.Getbitmap(0)
		//Desmonta el sistema de archivos
		lwh.UnmountVDisk()
		return m.CreateBitMap(b)
	case "bm_detdir":
		//Monta el sistema de archivos
		lwh.MountVDisk(m.diskPath, m.namePath)
		//Recupera el bitmap de detalle de directorio
		b := lwh.Getbitmap(1)
		//Desmonta el sistema de archivos
		lwh.UnmountVDisk()
		return m.CreateBitMap(b)
	case "bm_inode":
		//Monta el sistema de archivos
		lwh.MountVDisk(m.diskPath, m.namePath)
		//Recupera el bitmap de inodos
		b := lwh.Getbitmap(2)
		//Desmonta el sistema de archivos
		lwh.UnmountVDisk()
		return m.CreateBitMap(b)
	case "bm_block":
		//Monta el sistema de archivos
		lwh.MountVDisk(m.diskPath, m.namePath)
		//Recupera el bitmap de bloque de datos
		b := lwh.Getbitmap(3)
		//Desmonta el sistema de archivos
		lwh.UnmountVDisk()
		return m.CreateBitMap(b)
	case "bitacora":
		//Monta el sistema de archivos
		lwh.MountVDisk(m.diskPath, m.namePath)
		//Recupera el registro de la bitacora
		logArr := lwh.Getlogs()
		//Desmonta el sistema de archivos
		lwh.UnmountVDisk()
		//Crea el texto para el reporte
		return m.CreateLog(logArr)
	case "directorio":
		//Monta el sistema de archivos
		lwh.MountVDisk(m.diskPath, m.namePath)
		//Recupera el texto
		b := m.CreateDir()
		lwh.UnmountVDisk()
		return b
	}
	return []byte{}
}
