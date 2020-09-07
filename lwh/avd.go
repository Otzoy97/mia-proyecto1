package lwh

import "time"

//Avd ...
type Avd struct {
	FechaCreacion            [15]byte
	NombreDirectorio         [20]byte
	ApArraySubdirectorios    [6]int32
	ApDetalleDirectorio      int32
	ApArbolVirtualDirectorio int32
	Proper                   int32
	Auth                     [3]byte
}

//New configura un nuevo Arbol virtual de directorio
func (a *Avd) New(name, auth string, proper int) {
	//Establece la fecha de creación
	tm, _ := time.Now().GobEncode()
	copy(a.FechaCreacion[:], tm)
	copy(a.NombreDirectorio[:], name)
	a.ApArbolVirtualDirectorio = -1
	a.Proper = int32(proper)
	copy(a.Auth[:], auth)
}

//Extend ...
// func (a *Avd) Extend() Avd {
// 	ret := Avd{AvdNombreDirectorio: a.AvdNombreDirectorio, AvdProper: a.AvdProper, }
// }
