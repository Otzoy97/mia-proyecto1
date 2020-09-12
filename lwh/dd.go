package lwh

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

//Dd ...
type Dd struct {
	ArrayFiles          [5]DdFile
	ApDetalleDirectorio int32
}

//ReadDd recupera la información del detalle de directorio
//Par que tenga exito el puntero de archivo debe estar
//previamente puesto en posición con file.Seek
func (d *Dd) ReadDd() bool {
	arr := make([]byte, int(unsafe.Sizeof(*d)))
	if _, err := virtualDisk.Read(arr); err != nil {
		return false
	}
	buff := bytes.NewBuffer(arr)
	if err := binary.Read(buff, binary.BigEndian, d); err != nil {
		return false
	}
	return true
}

//tourDD busca en una coincidencia con el nombre dado
//devuelve un puntero al primer inodo del archivo
func (d *Dd) tourDD(name string) int64 {
	//Si el nombre es vacío retorna -1
	if name == "" {
		return -1
	}
	for _, vDet := range d.ArrayFiles {
		//recupera el nombre de vdet
		idxEnd := bytes.IndexByte(vDet.FileNombre[:], 0)
		if idxEnd == -1 {
			idxEnd = 25
		}
		tempName := string(vDet.FileNombre[:idxEnd])
		if tempName == name {
			return int64(vDet.FileApInodo)
		}
	}
	//Si el apDetalledirectorio no es cero
	if d.ApDetalleDirectorio != 0 {
		//Se mueve al puntero y realiza una búsqueda
		offset := int64(vdSuperBoot.SbApDetalleDirectorio) + int64(d.ApDetalleDirectorio)*int64(unsafe.Sizeof(Dd{}))
		virtualDisk.Seek(offset, 0)
		var dd Dd
		//Lee el detalle de directorio
		if !dd.ReadDd() {
			//No se pudo leer
			return -1
		}
		//Busca recursivamente el archivo
		return dd.tourDD(name)
	}
	//No encontró nada, devuelve -1
	return -1
}
