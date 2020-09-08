package lwh

import "time"

//Operaciones
const (
	MKDIR = iota
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
