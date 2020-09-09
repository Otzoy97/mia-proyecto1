package disk

import (
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"

	"github.com/fatih/color"
)

//Ebr ...
type Ebr struct {
	PartStatus, PartFit byte
	PartStart, PartSize uint32
	PartNext            int32
	PartName            [16]byte
}

//WriteEbr escribe el ebr en el archivo f en la posición que indica
func (e *Ebr) WriteEbr(f *os.File) bool {
	bin := new(bytes.Buffer)
	binary.Write(bin, binary.BigEndian, e)
	//Coloca el puntero en posición
	f.Seek(int64(e.PartStart), 0)
	if _, err := f.Write(bin.Bytes()); err != nil {
		color.New(color.FgHiYellow).Printf("No se pudo escribir el ebr %v\n", f.Name())
		return false
	}
	return true
}

//ReadEbr ...
func (e *Ebr) ReadEbr(f *os.File, startByte int) bool {
	dArr := make([]byte, int(unsafe.Sizeof(*e)))
	if _, err := f.Read(dArr); err != nil {
		color.New(color.FgHiYellow).Printf("No se pudo recuperar el ebr %v\n", f.Name())
		return false
	}
	buff := bytes.NewBuffer(dArr)
	if err := binary.Read(buff, binary.BigEndian, e); err != nil {
		color.New(color.FgHiYellow).Printf("No se pudo recuperar el ebr %v\n", f.Name())
		return false
	}
	return true
}
