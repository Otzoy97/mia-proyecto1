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

//Getlogs recupera todos los registros de log
func Getlogs() []Log {
	//Se posiciona al inicio del log y lee un registro de log
	pLog := vdSuperBoot.SbApLog
	virtualDisk.Seek(int64(pLog), 0)
	var lg Log
	//Le el primer Log
	//Siempre habrá al menos 2 registros del log
	lg.ReadLog()
	//Prepara el array en donde se almacenarán los datos
	var blog []Log
	//Realiza un ciclo leyendo todos los siguientes log
	//Si un log está "vacío" tendra LogTipoOperación = 0
	//lo cual no es posible ya que solo puede tomar operaciones
	//con enteros mayor a 0
	for lg.LogTipoOperacion > 0 {
		//Agrega el log que se leyó en la iteración anterior
		blog = append(blog, lg)
		//Lee un nuevo log
		lg.ReadLog()
	}
	return blog
}
