package fs

import (
	"bytes"
	"encoding/binary"
	"mia-proyecto1/disk"
	"mia-proyecto1/lwh"
	"os"
	"strings"

	"github.com/fatih/color"
)

//Mkfs ...
type Mkfs struct {
	id     string
	tipo   string
	extype byte
	Row    int
}

//AddOp añade un parámetro al comando Mkfs
func (m *Mkfs) AddOp(s string, v interface{}) {
	if s == "id" {
		m.id = v.(string)
	} else if s == "type" || s == "tipo" {
		m.tipo = v.(string)
	}
}

//Validate verifica que Mkfs tenga los parámetros necesarios
func (m *Mkfs) Validate() bool {
	if m.tipo == "" {
		m.extype = 'a'
	} else {
		if m.tipo == "full" {
			m.extype = 'u'
		} else if m.tipo == "fast" {
			m.extype = 'a'
		} else {
			color.New(color.FgHiYellow).Printf("Mkfs: type/tipo solo puede tomar el valor de fast y full (%v)\n", m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Mkfs no se puede ejecutar")
			return false
		}
	}
	if m.id == "" {
		color.New(color.FgHiYellow).Printf("Mkfs: no se encontró id (%v)\n", m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Mkfs no se puede ejecutar")
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
	//Configura y escribe las estructuras en el disco
	if !m.setStructs(file, name, path) {
		color.New(color.FgHiRed, color.Bold).Println("Mkdisk fracasó")
		return
	}
	//Cierra el archivo
	file.Close()
	color.New(color.FgHiGreen, color.Bold).Printf("Mkfs: el sistema de archivo en '%v' - '%v' se configuró correctamente\n", name, path)
}

//Crea y configura las estructuras para el sistema de archivos
func (m *Mkfs) setStructs(file *os.File, name string, path string) bool {
	//Recupera la información del mbr y la partición
	var mbr disk.Mbr
	mbr.ReadMbr(file)
	parArr, _ := mbr.CreateArrPart()
	//Se asume que la partición DEBE existir
	//Recupera la partición
	par := parArr.Find(name)
	//Limpia la partición
	buff := make([]byte, int(par.PartSize))
	file.Seek(0, int(par.PartStart))
	if _, err := file.Write(buff); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: el disco no se pudo formatear (%v)\n", m.Row)
		return false
	}
	//Realiza un split para recuperar el nombre virtual del disco
	strSplt := strings.Split(path, "/")
	hdName := strings.Split(strSplt[len(strSplt)-1], ".")[0]
	//Configura el superboot
	var sb lwh.Superboot
	sb.New(par, hdName)
	//Crea el AVD raíz
	var avd lwh.Avd
	avd.NewAvd("/", "664", 1, 1)
	//Crea un detalle de directorio para la raiz
	dd := lwh.Dd{ApDetalleDirectorio: -1}
	//Crea el archivo user.txt para el directorio raiz
	dd.ArrayFiles[0].NewDdFile("user.txt")
	//Crea los bloques de texto
	var db0 lwh.DataBlock
	var db1 lwh.DataBlock
	copy(db0.Data[:], "1,G,root\n1,U,root,root,2")
	copy(db1.Data[:], "01602782\n")
	//Crea los inodos
	var ino lwh.Inodo
	ino.NewInodo(0, 1, 1, "777")
	ino.ArrayBloques[0] = 0
	ino.ArrayBloques[1] = 1
	//Actualiza el superboot
	sb.SbInodosFree--
	sb.SbArbolVirtualFree--
	sb.SbDetalleDirectorioFree--
	sb.SbBloquesFree -= 2
	//Escribe el superboot al inicio de la partición
	bin := new(bytes.Buffer)
	file.Seek(0, int(par.PartStart))
	binary.Write(bin, binary.BigEndian, &sb)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el superboot (%v)\n", m.Row)
		return false
	}
	//Escribe un 1 en el bitmap de avd
	bin.Reset()
	var b byte = 1
	file.Seek(0, int(sb.SpApBitMapArbolDirectorio))
	binary.Write(bin, binary.BigEndian, b)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el bitmap del Avd (%v)\n", m.Row)
		return false
	}
	//Escribe el avd
	bin.Reset()
	file.Seek(0, int(sb.SbApArbolDirectorio))
	binary.Write(bin, binary.BigEndian, &avd)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el Avd (%v)\n", m.Row)
		return false
	}
	//Escribe un 1 en el bitmap de detalle de directorio
	bin.Reset()
	file.Seek(0, int(sb.SbApBitmapDetalleDirectorio))
	binary.Write(bin, binary.BigEndian, b)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el bitmap de detalle de directorio (%v)\n", m.Row)
		return false
	}
	//Escribe el detalle de directorio
	bin.Reset()
	file.Seek(0, int(sb.SbApDetalleDirectorio))
	binary.Write(bin, binary.BigEndian, &dd)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el detalle de directorio (%v)\n", m.Row)
		return false
	}
	//Escribe el bit map inodo
	bin.Reset()
	file.Seek(0, int(sb.SbApBitMapaTablaInodo))
	binary.Write(bin, binary.BigEndian, b)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el bitmap de inodo (%v)\n", m.Row)
		return false
	}
	//Escribe el inodo
	bin.Reset()
	file.Seek(0, int(sb.SbApTablaInodo))
	binary.Write(bin, binary.BigEndian, &ino)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el inodo (%v)\n", m.Row)
		return false
	}
	//Escribe el bit map de bloque de datos
	bin.Reset()
	file.Seek(0, int(sb.SbApBitmapBloques))
	binary.Write(bin, binary.BigEndian, b)
	binary.Write(bin, binary.BigEndian, b)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el bitmap de bloque de datos (%v)\n", m.Row)
		return false
	}
	//Escribe el bloque de datos
	bin.Reset()
	file.Seek(0, int(sb.SbApBloques))
	binary.Write(bin, binary.BigEndian, &db0)
	binary.Write(bin, binary.BigEndian, &db1)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el bloque de datos (%v)\n", m.Row)
		return false
	}
	//Escribe el LOG
	var log01 lwh.Log
	var log02 lwh.Log
	log01.NewLog(lwh.MKDIR, 0, "/", "")
	log02.NewLog(lwh.MKDIR, 0, "/user.txt", "1,G,root\n1,U,root,root,201602782")
	bin.Reset()
	file.Seek(0, int(sb.SbApLog))
	binary.Write(bin, binary.BigEndian, &log01)
	binary.Write(bin, binary.BigEndian, &log02)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir el inodo (%v)\n", m.Row)
		return false
	}
	//Escriba la copia del superboot
	bin.Reset()
	file.Seek(0, int(sb.SbApLog+sb.SbArbolvirtualCount*sb.SbApLog))
	binary.Write(bin, binary.BigEndian, &sb)
	if _, err := file.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("Mkfs: no se pudo escribir la copia del superboot (%v)\n", m.Row)
		return false
	}
	//Informa al usuario del resultado
	return true
}
