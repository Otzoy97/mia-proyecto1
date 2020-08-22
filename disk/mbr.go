package disk

import (
	"bytes"
	"encoding/binary"
	"os"
	"time"
	"unsafe"

	"github.com/fatih/color"
)

//Mbr ...
type Mbr struct {
	MbrTamanio, MbrDiskSignature int
	MbrFechaCreacion             time.Time
	MbrPartition1                Partition
	MbrPartition2                Partition
	MbrPartition3                Partition
	MbrPartition4                Partition
}

//WriteMbr escribe el mbr en el archivo f
func (m *Mbr) WriteMbr(f *os.File) bool {
	var bin bytes.Buffer
	binary.Write(&bin, binary.BigEndian, m)
	f.Seek(0, 0)
	if _, err := f.WriteAt(bin.Bytes(), 0); err != nil {
		color.New(color.FgHiYellow).Printf("No se pudo escribir el mbr %v\n", f.Name())
		return false
	}
	return true
}

//ReadMbr lee el mbr del archivo f
func (m *Mbr) ReadMbr(f *os.File) bool {
	dArr := make([]byte, int(unsafe.Sizeof(*m)))
	if _, err := f.Read(dArr); err != nil {
		color.New(color.FgHiYellow).Printf("No se pudo recuperar el mbr %v\n", f.Name())
		return false
	}
	buff := bytes.NewBuffer(dArr)
	if err := binary.Read(buff, binary.BigEndian, m); err != nil {
		color.New(color.FgHiYellow).Printf("No se pudo recuperar el mbr %v\n", f.Name())
		return false
	}
	return true
}
