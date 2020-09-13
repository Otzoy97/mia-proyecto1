package lwh

import "github.com/fatih/color"

//Bitmap ...
type Bitmap []byte

//BmType ...
type BmType byte

//BmType ...
const (
	BitmapAvd   BmType = 1
	BitmapDd           = 2
	BitmapInodo        = 3
	BitmapBd           = 4
)

//Getbitmap recupera el stream de bytes que representan los espacios
//llenos y vacios del respectivo tipo de bitmap
func Getbitmap(op BmType) []byte {
	//Recupera el puntero del bitmap y el tamaño
	point, size := whichBM(op)
	//Coloca el punteo de disco en posición
	virtualDisk.Seek(int64(point), 0)
	b := make([]byte, size)
	//Lee el bit map
	if _, err := virtualDisk.Read(b); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el bitmap de directorio %v\n     %v\n", virtualDisk.Name(), err.Error())
	}
	return b
}

//whichBM decide el puntero a utilizar y la cantidad de bytes a leer
func whichBM(t BmType) (int32, int32) {
	switch t {
	case BitmapAvd:
		return vdSuperBoot.SbApBitMapArbolDirectorio, vdSuperBoot.SbArbolvirtualCount
	case BitmapBd:
		return vdSuperBoot.SbApBitmapBloques, vdSuperBoot.SbBloquesCount
	case BitmapDd:
		return vdSuperBoot.SbApBitmapDetalleDirectorio, vdSuperBoot.SbDetalleDirectorioCount
	case BitmapInodo:
		return vdSuperBoot.SbApBitMapaTablaInodo, vdSuperBoot.SbInodosCount
	}
	return 0, -1
}
