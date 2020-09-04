package lwh

//Avd ...
type Avd struct {
	AvdFechaCreacion            [15]byte
	AvdNombreDirectorio         [20]byte
	AvdApArraySubdirectorios    [6]int32
	AvdApDetalleDirectorio      int32
	AvdApArbolVirtualDirectorio int32
	AvdProper                   int32
}
