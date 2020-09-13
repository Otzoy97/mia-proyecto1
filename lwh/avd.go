package lwh

import (
	"bytes"
	"encoding/binary"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/fatih/color"
)

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
func (a *Avd) NewAvd(name, auth string, proper, gid int32) {
	//Establece la fecha de creación
	tm, _ := time.Now().GobEncode()
	copy(a.FechaCreacion[:], tm)
	copy(a.NombreDirectorio[:], name)
	a.ApArbolVirtualDirectorio = -1
	a.ApDetalleDirectorio = -1
	a.Proper = proper
	a.Gid = gid
	copy(a.Auth[:], auth)
}

//ReadAvd recupera la información del avd especificado
func (a *Avd) ReadAvd(n int32) bool {
	//Se mueve a la posición especificada
	offset := int64(vdSuperBoot.SbApArbolDirectorio + n*int32(unsafe.Sizeof(*a)))
	virtualDisk.Seek(offset, 0)
	//Lee el struct
	arr := make([]byte, int(unsafe.Sizeof(*a)))
	if _, err := virtualDisk.Read(arr); err != nil {
		return false
	}
	buff := bytes.NewBuffer(arr)
	if err := binary.Read(buff, binary.BigEndian, a); err != nil {
		return false
	}
	return true
}

//WriteAvd escribe el avd en al posición especificada
func (a *Avd) WriteAvd(n int32) bool {
	//Se mueve a la posición del disco
	offset := int64(vdSuperBoot.SbApArbolDirectorio + n*int32(unsafe.Sizeof(*a)))
	virtualDisk.Seek(offset, 0)
	//Escribe el struct en un stream de bytes
	bin := new(bytes.Buffer)
	binary.Write(bin, binary.BigEndian, a)
	//Escribe el struct en el disco
	if _, err := virtualDisk.Write(bin.Bytes()); err != nil {
		return false
	}
	return true
}

//Find busca recursivamente un archivo según el path dado.
//Devuelve el puntero y el tipo de dato (avd o inodo)
//0 - directorio
//1 - archivo
//2 - no se encontró
//3 - no es directorio
//4 - no se pudo leer algo
func (a *Avd) Find(path string) (int32, byte) {
	//Verifica si el path es igual a "/"
	if path == "/" {
		return -1, 2
	}
	//Valida el path
	nPath, flag := validatePath(path)
	if !flag {
		color.New(color.FgHiYellow).Printf("     '%v' no es un directorio\n", path)
		return -1, 3
	}
	//Recorre el avd buscando coincidencias
	return a.tourAVD(nPath)
}

//tourAVD busca en los puntero de avd y dd aluguna coincidencia con los nombres en
//el slice dir. Devuelve el puntero y el tipo de dato (avd o inodo)
//0 - directorio
//1 - archivo
//2 - no se encontró
//3 - no es directorio
//4 - no se pudo leer algo
func (a *Avd) tourAVD(dir []string) (int32, byte) {
	//No hay nombre para leer
	if len(dir) == 0 {
		return -1, 2
	}
	//Busca en el detalle de directorio
	if a.ApDetalleDirectorio != -1 {
		var dd Dd
		//Lee el detalle de directorio
		if dd.ReadDd(a.ApDetalleDirectorio) {
			//Busca el primer elemento del array
			off := dd.tourDD(dir[0])
			if off != -1 {
				return off, 1
			}
		} else {
			return -1, 4
		}
	}
	//Busca en el array del avd actual
	for _, pAvd := range a.ApArraySubdirectorios {
		//El avd existe
		if pAvd != 0 {
			//Si no es cero, entonces debe apuntar a algún lado
			var newAvd Avd
			//Lee el avd
			if !newAvd.ReadAvd(pAvd) {
				//No se leyó correctamente
				return -1, 4
			}
			//Convierte el array de bytes en string
			idxEnd := bytes.IndexByte(newAvd.NombreDirectorio[:], 0)
			if idxEnd == -1 {
				idxEnd = 20
			}
			tempName := string(newAvd.NombreDirectorio[:idxEnd])
			//Compara si el nombre es el mismo
			if tempName == dir[0] {
				//El primer nombre es igual
				if len(dir) > 1 {
					//Aun hay nombres por buscar
					return newAvd.tourAVD(dir[1:])
				}
				//Ya no hay nombres por buscar
				return pAvd, 0
			}
		}
	}
	//Busca en el apuntador indirecto
	if a.ApArbolVirtualDirectorio != -1 {
		//Mueve el apuntador de disco a esa posición y recupera el avd
		var newAvd Avd
		//Lee el avd
		if !newAvd.ReadAvd(a.ApArbolVirtualDirectorio) {
			//No se leyó correctamente
			return -1, 4
		}
		//Vuelve a llamar esta misma función
		return newAvd.tourAVD(dir)
	}
	//No se encontró ninguna coincidencia
	return -1, 2
}

//CreateDir true si logra crear el directorio, falso si no
//0 - ya existe el directorio
//1 - se logró crear el directorio
//2 - no se pudo crear el directorio
//3 - el path no era válido
//4 - no se pudo leer el disco
func (a *Avd) CreateDir(from string, to []string) (bool, int) {
	//Verifica que 'to', tenga datos
	if len(to) == 0 {
		//Si llega a este punto, es porque no se creó nada
		//0 - ya existe el directorio
		return false, 0
	}
	//Arma el path a buscar
	from += strings.Join(to[:1], "/")
	//Busca el path, utilizando 'from'
	pAvd, tipe := a.Find(from)
	if tipe == 3 {
		//3 - el path no era válido
		return false, 3
	} else if tipe == 4 {
		//4 - no se pudo leer el disco
		return false, 4
	} else if tipe == 0 {
		//Si encuentra algo debe ser un directorio, si no es un error
		var newAvd Avd
		//Recupera la información del nuevo avd
		newAvd.ReadAvd(pAvd)
		//Crea el directorio en ese avd
		return newAvd.CreateDir(from, to[1:])
	} else if tipe == 2 {
		//No se encontró, entonces debe crearse
		//Recupera el bitmap para el arbol de directorio
		avdBm := Bitmap(Getbitmap(BitmapAvd))
		if bmAp, flag := avdBm.FindSpaces(1); flag {
			//Crea un nuevo avd
			var newAvd Avd
			newAvd.NewAvd(to[0], "664", logUser.uid, logUser.gid)
			//Escribe el nuevo avd
			if newAvd.WriteAvd(bmAp) {
				//Escribe el bmp
				if writeBM(BitmapAvd, bmAp) {
					//Se creó correctamente el directorio
					//Verifica si to, aún tiene datos
					if len(to) > 0 {
						return newAvd.CreateDir(from, to[1:])
					}
					return true, 1
				}
				return false, 2
			}
			return false, 2
		}
	}
	return false, 2
}

//validatePath el offset es la posición la cual hay que leer
func validatePath(path string) ([]string, bool) {
	//Se asegura que sea un directorio válido
	match, _ := regexp.Match(`^(/[^/]*)+/?$`, []byte(path))
	if !match {
		//No es un directorio válido
		return []string{}, false
	}
	//Separa el directorio
	psplit := strings.Split(path, "/")
	//Quita el primer directorio
	return psplit[1:], true
}
