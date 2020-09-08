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
	Gid                      int32
}

//NewAvd configura un nuevo Arbol virtual de directorio
func (a *Avd) NewAvd(name, auth string, proper, gid int) {
	//Establece la fecha de creación
	tm, _ := time.Now().GobEncode()
	copy(a.FechaCreacion[:], tm)
	copy(a.NombreDirectorio[:], name)
	a.ApArbolVirtualDirectorio = -1
	a.Proper = int32(proper)
	a.Gid = int32(gid)
	copy(a.Auth[:], auth)
}

//Extend ...
// func (a *Avd) Extend() Avd {
// 	ret := Avd{AvdNombreDirectorio: a.AvdNombreDirectorio, AvdProper: a.AvdProper, }
// }
