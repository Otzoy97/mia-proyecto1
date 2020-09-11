package lwh

import "bytes"

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
