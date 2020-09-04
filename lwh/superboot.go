package lwh

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
