package lwh

import (
	"bytes"
	"encoding/binary"
	"time"
	"unsafe"

	"github.com/fatih/color"
)

//Operaciones
const (
	MKDIR = iota + 1
	MKFILE
)

//Log 0- carpeta// 1- archivo
type Log struct {
	LogTipoOperacion byte
	LogTipo          byte
	LogNombre        [256]byte
	LogContenido     [256]byte
	LogFecha         [15]byte
}

//NewLog ...
func (l *Log) NewLog(op, tipo byte, name, cont string) {
	l.LogTipoOperacion = op
	l.LogTipo = tipo
	copy(l.LogNombre[:], name)
	copy(l.LogContenido[:], cont)
	tim, _ := time.Now().GobEncode()
	copy(l.LogFecha[:], tim)
}

//ReadLog lee un registro de log desde virtualDisk
//Antes de leer debe colocarse el puntero de archivo
//en la posición del ap_log
func (l *Log) ReadLog() bool {
	//Crea un arreglo de bytes del tamaño del struct de lob
	darr := make([]byte, int(unsafe.Sizeof(*l)))
	if _, err := virtualDisk.Read(darr); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el registro de Log %v\n", virtualDisk.Name())
		return false
	}
	bin := bytes.NewBuffer(darr)
	if err := binary.Read(bin, binary.BigEndian, l); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el registro de Log %v\n", virtualDisk.Name())
		return false
	}
	return true
}

//Getop devuelve el nombre de la operación realizada
func (l *Log) Getop() string {
	switch l.LogTipoOperacion {
	case 1:
		return "MKDIR"
	case 2:
		return "MKFILE"
	}
	return "--"
}

//Getipo devuelve el tipo de objeto creado (carpeta/archivo)
func (l *Log) Getipo() string {
	switch l.LogTipo {
	case 0:
		return "DIR"
	case 1:
		return "FILE"
	}
	return "--"
}
