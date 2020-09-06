package disk

import (
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"

	"github.com/fatih/color"
)

//Mbr ...
type Mbr struct {
	MbrTamanio, MbrDiskSignature uint32
	MbrFechaCreacion             [15]byte
	MbrPartition1                Partition
	MbrPartition2                Partition
	MbrPartition3                Partition
	MbrPartition4                Partition
}

//WriteMbr escribe el mbr en el archivo f
func (m *Mbr) WriteMbr(f *os.File) bool {
	bin := new(bytes.Buffer)
	binary.Write(bin, binary.BigEndian, m)
	f.Seek(0, 0)
	if _, err := f.WriteAt(bin.Bytes(), 0); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo escribir el mbr %v\n", f.Name())
		return false
	}
	return true
}

//ReadMbr lee el mbr del archivo f
func (m *Mbr) ReadMbr(f *os.File) bool {
	//Crea un arreglo de byes del tamaño del struct del mbr
	dArr := make([]byte, int(unsafe.Sizeof(*m)))
	if _, err := f.Read(dArr); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el mbr %v\n", f.Name())
		return false
	}
	buff := bytes.NewBuffer(dArr)
	if err := binary.Read(buff, binary.BigEndian, m); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el mbr %v\n", f.Name())
		return false
	}
	return true
}

//CreateArrPart crea una array de particiones con las particiones
//que tienen status = 1, devuelve true si hay más de una particion
//extendida
func (m *Mbr) CreateArrPart() (ByPartStart, bool) {
	flag := false
	var arrPar ByPartStart = []Partition{}
	if m.MbrPartition1.PartStatus == 1 {
		arrPar = append(arrPar, m.MbrPartition1)
		flag = (flag || m.MbrPartition1.PartType == 'e')
	}
	if m.MbrPartition2.PartStatus == 1 {
		arrPar = append(arrPar, m.MbrPartition2)
		flag = (flag || m.MbrPartition2.PartType == 'e')
	}
	if m.MbrPartition3.PartStatus == 1 {
		arrPar = append(arrPar, m.MbrPartition3)
		flag = (flag || m.MbrPartition3.PartType == 'e')
	}
	if m.MbrPartition4.PartStatus == 1 {
		arrPar = append(arrPar, m.MbrPartition4)
		flag = (flag || m.MbrPartition4.PartType == 'e')
	}
	return arrPar, flag
}
