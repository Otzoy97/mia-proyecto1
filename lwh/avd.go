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

//ReadAvd ...
//Para que la lectura sea exitosa, el apuntador de virtualdisk
//debe estar en la posición de un avd
func (a *Avd) ReadAvd() bool {
	arr := make([]byte, int(unsafe.Sizeof(*a)))
	if _, err := virtualDisk.Read(arr); err != nil {
		return false
	}
	buff := bytes.NewBuffer(arr)
	if err := binary.Read(buff, binary.BigEndian, a); err != nill {
		return false
	}
	return true
}

//Find busca recursivamente un archivo según el path dado
func (a *Avd) Find(path string) int32 {
	//Verifica si el path es igual a "/"
	if path == "/" {
		return 0
	}
	//Valida el path
	nPath, flag := validatePath(path)
	if !flag {
		color.New(color.FgHiYellow).Printf("     '%v' no es un directorio\n", path)
		return -1
	}

	//Recorre el avd buscando coincidencias
}

//Tour los puntero de avd, buscando coincidencia con los nombre en
//dir. Devuelve el punteo al último avd al que hace coincidencia
func (a *Avd) Tour(dir []string) int64 {
	//No hay nombre para leer
	if len(dir) == 0 {
		return -1
	}
	//Busca en el array del avd actual
	for _, avd := range a.ApArraySubdirectorios {
		//El avd existe
		if avd != 0 {
			//Si no es cero, entonces debe apuntar a algún lado
			//Mueve el apuntador de disco a esa posición y recupera el avd
			offset := int64(vdSuperBoot.SbApArbolDirectorio + avd*int32(unsafe.Sizeof(Avd{})))
			virtualDisk.Seek(offset, 0)
			var newAvd Avd
			//Lee el avd
			if !newAvd.ReadAvd() {
				//No se leyó correctamente
				return -1
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
					return newAvd.Tour(dir[1:])
				}
				//Ya no hay nombres por buscar
				return offset
			}
		}
	}
	//Busca en el apuntador indirecto
	if a.ApArbolVirtualDirectorio != -1 {
		//Mueve el apuntador de disco a esa posición y recupera el avd
		offset := int64(vdSuperBoot.SbApArbolDirectorio + a.ApArbolVirtualDirectorio*int32(unsafe.Sizeof(Avd{})))
		virtualDisk.Seek(offset, 0)
		var newAvd Avd
		//Lee el avd
		if !newAvd.ReadAvd() {
			//No se leyó correctamente
			return -1
		}
		//Vuelve a llamar esta misma función
		return newAvd.Tour(dir)
	}
	//No se encontró ninguna coincidencia
	return -1
}

//validatePath el offset es la posición la cual hay que leer
//el tipo devuelve si es una avd o un inodo
func validatePath(path string) ([]string, bool) {
	//Se asegura que sea un directorio válido
	match, _ := regexp.Match(`^(/[^/ ]*)+/?$`, []byte(path))
	if !match {
		//No es un directorio válido
		return []string{}, false
	}
	//Separa el directorio
	psplit := strings.Split(path, "/")
	//Quita el primer directorio
	return psplit[1:], true
}
