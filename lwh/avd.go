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
	a.ApDetalleDirectorio = -1
	a.Proper = int32(proper)
	a.Gid = int32(gid)
	copy(a.Auth[:], auth)
}

//ReadAvd recupera la información del avd especificado
func (a *Avd) ReadAvd(n int32) (bool, int64) {
	//Se mueve a la posición especificada
	offset := int64(vdSuperBoot.SbApArbolDirectorio + n*int32(unsafe.Sizeof(*a)))
	virtualDisk.Seek(offset, 0)
	//Lee el struct
	arr := make([]byte, int(unsafe.Sizeof(*a)))
	if _, err := virtualDisk.Read(arr); err != nil {
		return false, -1
	}
	buff := bytes.NewBuffer(arr)
	if err := binary.Read(buff, binary.BigEndian, a); err != nil {
		return false, -1
	}
	return true, offset
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

//Find busca recursivamente un archivo según el path dado
func (a *Avd) Find(path string) (int64, byte) {
	//Verifica si el path es igual a "/"
	if path == "/" {
		return -1, 2
	}
	//Valida el path
	nPath, flag := validatePath(path)
	if !flag {
		color.New(color.FgHiYellow).Printf("     '%v' no es un directorio\n", path)
		return -1, 2
	}
	//Recorre el avd buscando coincidencias
	return a.tourAVD(nPath)
}

//tourAVD busca en los puntero de avd y dd aluguna coincidencia con los nombres en
//el slice dir. Devuelve el puntero y el tipo de dato (avd o dd)
//0 - directorio 1 - archivo  2 - error
func (a *Avd) tourAVD(dir []string) (int64, byte) {
	//No hay nombre para leer
	if len(dir) == 0 {
		return -1, 2
	}
	//Busca en el detalle de directorio
	if a.ApDetalleDirectorio != -1 {
		//Se coloca en la posición del apuntador de detalle de directorio
		offset := int64(vdSuperBoot.SbApArbolDirectorio + a.ApDetalleDirectorio*int32(unsafe.Sizeof(Dd{})))
		virtualDisk.Seek(offset, 0)
		var dd Dd
		//Lee el detalle de directorio
		dd.ReadDd()
		//Busca el primer elemento del array
		off := dd.tourDD(dir[0])
		if off != -1 {
			return off, 1
		}
	}
	//Busca en el array del avd actual
	for _, avd := range a.ApArraySubdirectorios {
		//El avd existe
		if avd != 0 {
			//Si no es cero, entonces debe apuntar a algún lado
			var newAvd Avd
			//Lee el avd
			flag, offset := newAvd.ReadAvd(avd)
			if !flag {
				//No se leyó correctamente
				return -1, 2
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
				return offset, 0
			}
		}
	}
	//Busca en el apuntador indirecto
	if a.ApArbolVirtualDirectorio != -1 {
		//Mueve el apuntador de disco a esa posición y recupera el avd
		var newAvd Avd
		//Lee el avd
		flag, _ := newAvd.ReadAvd(a.ApArbolVirtualDirectorio)
		if !flag {
			//No se leyó correctamente
			return -1, 2
		}
		//Vuelve a llamar esta misma función
		return newAvd.tourAVD(dir)
	}
	//No se encontró ninguna coincidencia
	return 0, 2
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
