package lwh

import (
	"bytes"
	"time"
)

//Dd ...
type Dd struct {
	ArrayFiles          [5]DdFile
	ApDetalleDirectorio int32
}

//DdFile ...
type DdFile struct {
	FileNombre          [25]byte
	FileApInodo         int32
	FileDateCreacion    [15]byte
	FileDateModficacion [15]byte
}

//DataBlock ...
type DataBlock struct {
	Data [25]byte
}

//New configura los atributos de DdFile
func (d *DdFile) New(name string) {
	copy(d.FileNombre[:], name)
	tm, _ := time.Now().GobEncode()
	copy(d.FileDateCreacion[:], tm)
	copy(d.FileDateModficacion[:], tm)
}

//getBdData convierte el contenido del archivo en una cadena
func (d *DataBlock) getBdData() string {
	//Recupera los bytes hasta encontrar un caracter nulo
	idxEnd := bytes.IndexByte(d.Data[:], 0)
	if idxEnd == -1 {
		idxEnd = len(d.Data)
	}
	temName := d.Data[:idxEnd]
	return string(temName)
}
