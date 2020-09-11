package lwh

import "time"

//DdFile ...
type DdFile struct {
	FileNombre          [25]byte
	FileApInodo         int32
	FileDateCreacion    [15]byte
	FileDateModficacion [15]byte
}

//NewDdFile configura los atributos de DdFile
func (d *DdFile) NewDdFile(name string) {
	copy(d.FileNombre[:], name)
	tm, _ := time.Now().GobEncode()
	copy(d.FileDateCreacion[:], tm)
	copy(d.FileDateModficacion[:], tm)
}
