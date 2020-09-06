package lwh

import (
	"mia-proyecto1/disk"
	"time"
	"unsafe"
)

//Superboot ...
type Superboot struct {
	SbNombreHd                      [16]byte
	SbArbolvirtualCount             int32
	SbDetalleDirectorioCount        int32
	SbInodosCount                   int32
	SbBloquesCount                  int32
	SbArbolVirtualFree              int32
	SbDetalleDirectorioFree         int32
	SbInodosFree                    int32
	SbBloquesFree                   int32
	SbDateCreacion                  [15]byte
	SbDateUltimoMontaje             [15]byte
	SbMontajesCount                 int32
	SpApBitMapArbolDirectorio       int32
	SbApArbolDirectorio             int32
	SbApBitmapDetalleDirectorio     int32
	SbApDetalleDirectorio           int32
	SbApBitMapaTablaInodo           int32
	SbApTablaInodo                  int32
	SbApBitmapBloques               int32
	SbApBloques                     int32
	SbApLog                         int32
	SbSizeStructArbolDirectorio     int32
	SbSizeStructDetalleDirectorio   int32
	SbSizeStructInodo               int32
	SbSizeStructBloque              int32
	SbFirstFreeBitArbolDirectorio   int32
	SbFirstFreeBitDetalleDirectorio int32
	SbFirstFreeBitTablaInodo        int32
	SbFirstFreeBitBloques           int32
	SbMagicNum                      int32
}

//New calcula el tamaño de las estructuras de datos
//con el tamaño 'size' del disco. Name es el nombre del disco
func (s *Superboot) New(part disk.Partition, name string) {

	copy(s.SbNombreHd[:], name)
	size := int(part.PartSize)
	n := int32((size - (2 * int(unsafe.Sizeof(Superboot{})))) / (27 + int(unsafe.Sizeof(Avd{})) + int(unsafe.Sizeof(Dd{})) + 5*int(unsafe.Sizeof(Inodo{})) + 20*int(unsafe.Sizeof(DataBlock{})) + int(unsafe.Sizeof(Log{}))))
	s.SbArbolvirtualCount = n
	s.SbDetalleDirectorioCount = n
	s.SbInodosCount = 5 * n
	s.SbBloquesCount = 20 * n
	s.SbArbolVirtualFree = n
	s.SbDetalleDirectorioFree = n
	s.SbInodosFree = 5 * n
	s.SbBloquesFree = 20 * n
	//La fecha de ulitmo montaje conincidira
	//por esta única vez con la fecha de creación
	tDec, _ := time.Now().GobEncode()
	copy(s.SbDateCreacion[:], tDec)
	copy(s.SbDateUltimoMontaje[:], tDec)
	s.SbMontajesCount = 0
	s.SpApBitMapArbolDirectorio = int32(part.PartStart) + int32(unsafe.Sizeof(Superboot{}))
	s.SbApArbolDirectorio = s.SpApBitMapArbolDirectorio + n
	s.SbApBitmapDetalleDirectorio = s.SbApArbolDirectorio + n*int32(unsafe.Sizeof(Avd{}))
	s.SbApDetalleDirectorio = s.SbApBitmapDetalleDirectorio + n
	s.SbApBitMapaTablaInodo = s.SbApDetalleDirectorio + n*int32(unsafe.Sizeof(Dd{}))
	s.SbApTablaInodo = s.SbApBitMapaTablaInodo + 5*n
	s.SbApBitmapBloques = s.SbApTablaInodo + 5*n*int32(unsafe.Sizeof(Inodo{}))
	s.SbApBloques = s.SbApBitmapBloques + 20*n
	s.SbApLog = s.SbApBloques + 20*n*int32(unsafe.Sizeof(DataBlock{}))
	s.SbSizeStructArbolDirectorio = int32(unsafe.Sizeof(Avd{}))
	s.SbSizeStructDetalleDirectorio = int32(unsafe.Sizeof(Dd{}))
	s.SbSizeStructInodo = int32(unsafe.Sizeof(Inodo{}))
	s.SbSizeStructBloque = int32(unsafe.Sizeof(DataBlock{}))
	s.SbFirstFreeBitArbolDirectorio = 0
	s.SbFirstFreeBitDetalleDirectorio = 0
	s.SbFirstFreeBitTablaInodo = 0
	s.SbFirstFreeBitBloques = 0
	s.SbMagicNum = 201602782

}