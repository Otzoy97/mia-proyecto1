package lwh

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/fatih/color"
)

//Dd ...
type Dd struct {
	ArrayFiles          [5]DdFile
	ApDetalleDirectorio int32
}

//ReadDd recupera la información del detalle de directorio
//Par que tenga exito el puntero de archivo debe estar
//previamente puesto en posición con file.Seek
func (d *Dd) ReadDd(n int32) bool {
	//Se coloca en la posición del apuntador de detalle de directorio
	offset := int64(vdSuperBoot.SbApDetalleDirectorio + n*int32(unsafe.Sizeof(Dd{})))
	virtualDisk.Seek(offset, 0)
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

//WriteDd ...
func (d *Dd) WriteDd(n int32) bool {
	//Se coloca en posición del aputnado
	offset := int64(vdSuperBoot.SbApDetalleDirectorio + n*int32(unsafe.Sizeof(*d)))
	virtualDisk.Seek(offset, 0)
	//Escribe el struct en un stream de bytes
	bin := new(bytes.Buffer)
	binary.Write(bin, binary.BigEndian, d)
	//Escrib el estruct en e disco
	if _, err := virtualDisk.Write(bin.Bytes()); err != nil {
		return false
	}
	return true

}

//CreateFile ...
func (d *Dd) CreateFile(name, cont string, size int32) {
	//Verifica si hay espacio para un nuevo archivo

	for i := 0; i < 5; i++ {
		//Para determinar si el dfil está ocupado
		//Se comparará el indice del elemento 0 para filenombre
		idxEnd := bytes.IndexByte(d.ArrayFiles[i].FileNombre[:], 0)
		if idxEnd == 0 {
			//Dfile disponible
			// d.ArrayFiles[i].NewDdFile(name)
			// d.ArrayFiles[i].
		}
	}
	//Utiliza el apuntador indirecto
	if d.ApDetalleDirectorio == -1 {
		//Busca espacio
		bm := Getbitmap(BitmapDd)
		pDD, flag := bm.FindSpaces(1)
		if !flag {
			color.New(color.FgHiRed, color.Bold).Println("Mkfile: no hay espacio disponible")
			return
		}
		//Escribe el nuevo detalle de directorio
		nuevoDd := Dd{ApDetalleDirectorio: -1}
		nuevoDd.ReadDd(pDD)
	}
	//Lee el apuntador indirecto
	var newDd Dd
	newDd.ReadDd(d.ApDetalleDirectorio)
	newDd.CreateFile(name, cont, size)
}

//tourDD busca en una coincidencia con el nombre dado
//devuelve un puntero al primer inodo del archivo
func (d *Dd) tourDD(name string) int32 {
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
			return vDet.FileApInodo
		}
	}
	//Si el apDetalledirectorio no es cero
	if d.ApDetalleDirectorio != 0 {
		var dd Dd
		//Lee el detalle de directorio
		if !dd.ReadDd(d.ApDetalleDirectorio) {
			//No se pudo leer
			return -1
		}
		//Busca recursivamente el archivo
		return dd.tourDD(name)
	}
	//No encontró nada, devuelve -1
	return -1
}
