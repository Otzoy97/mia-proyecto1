package lwh

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unsafe"
)

//Inodo ...
type Inodo struct {
	CountInodo            int32
	SizeArchivo           int32
	CountBloquesAsignados int32
	ArrayBloques          [4]int32
	ApIndirecto           int32
	IDProper              int32
	Auth                  [3]byte
	Gid                   int32
}

//NewInodo crea un nuevo inodo
func (i *Inodo) NewInodo(pInodo, proper, gid int, auth string) {
	i.CountInodo = int32(pInodo)
	i.SizeArchivo = 0
	i.CountBloquesAsignados = 0
	i.ArrayBloques = [4]int32{-1, -1, -1, -1}
	i.ApIndirecto = -1
	i.IDProper = int32(proper)
	i.Gid = int32(gid)
	copy(i.Auth[:], auth)
}

//getCont lee cada uno de los bloques y concatena el
//contenido que alojan
func (i *Inodo) getCont() string {
	//Alojará el contenido del archivo
	var strd strings.Builder
	//Viaja a la posición del inodo que especifica el
	//array de IArrayBloques y recupera el contenido binario
	for _, offset := range i.ArrayBloques {
		if offset != -1 {
			//El bloque está asignado
			var block DataBlock
			if !block.readDB(offset) {
				//No se pudo leer
				continue
			}
			//Añade el contenido del archivo a la cadena a retornar
			strd.WriteString(block.getBdData())
		}
	}
	//Verifica si el apuntador indirecto, está a puntando a algo
	if i.ApIndirecto != -1 {
		//Lee el inodo
		var in Inodo
		if in.readInodo(i.ApIndirecto) {
			//Concatena el resultado
			strd.WriteString(in.getCont())
		}
	}
	return strd.String()
}

//readInodo lee el inodo en la posición n
func (i *Inodo) readInodo(n int32) bool {
	//Se mueve a la posición especificada
	offset := int64(vdSuperBoot.SbApTablaInodo + n*int32(unsafe.Sizeof(*i)))
	virtualDisk.Seek(offset, 0)
	//Lee el struct
	arr := make([]byte, int(unsafe.Sizeof(*i)))
	if _, err := virtualDisk.Read(arr); err != nil {
		return false
	}
	buff := bytes.NewBuffer(arr)
	if err := binary.Read(buff, binary.BigEndian, i); err != nil {
		return false
	}
	return true
}

//writeIniodo lee el inodo en la posición n
func (i *Inodo) writeInodo(n int32) bool {
	//Se mueve a la posicion del disco
	offset := int64(vdSuperBoot.SbApTablaInodo + n*int32(unsafe.Sizeof(*i)))
	virtualDisk.Seek(offset, 0)
	//Lee le struct
	bin := new(bytes.Buffer)
	binary.Write(bin, binary.BigEndian, i)
	//Escrib el struct en el disco
	if _, err := virtualDisk.Write(bin.Bytes()); err != nil {
		return false
	}
	return true
}
