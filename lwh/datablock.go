package lwh

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

//DataBlock ...
type DataBlock struct {
	Data [25]byte
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

//readDB lee el struct Datablock desde el archivo virtualDisk
func (d *DataBlock) readDB(n int32) bool {
	offset := int64(vdSuperBoot.SbApBloques + n*int32(unsafe.Sizeof(*d)))
	virtualDisk.Seek(offset, 0)
	darr := make([]byte, int(unsafe.Sizeof(*d)))
	if _, err := virtualDisk.Read(darr); err != nil {
		return false
	}
	buff := bytes.NewBuffer(darr)
	if err := binary.Read(buff, binary.BigEndian, d); err != nil {
		return false
	}
	return true
}

//writeDB escribe el struct a el archivo virtualDisk
func (d *DataBlock) writeDB(n int32) bool {
	//Se mueve a la posición del disco
	offset := int64(vdSuperBoot.SbApBloques + n*int32(unsafe.Sizeof(*d)))
	virtualDisk.Seek(offset, 0)
	//Escribe el struct en un stream de bytes
	bin := new(bytes.Buffer)
	binary.Write(bin, binary.BigEndian, d)
	//Escribe el struct en el disco
	if _, err := virtualDisk.Write(bin.Bytes()); err != nil {
		return false
	}
	return true
}
