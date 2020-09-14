package lwh

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

//CreateDir ...
func (a *Avd) CreateDir(from string, to []string, pos int32) {
	//Verifica que to aún tenga datos
	if len(to) == 0 {
		return
	}
	//Arma el path a buscar
	from += strings.Join(to[:1], "/")
	//Busca el path
	pAvd, tipe := a.Find(from)
	if tipe == 0 {
		//El directorio existe
		var avd Avd
		//Lee el apuntador del directorio
		avd.ReadAvd(pAvd)
		//Crea el siguiente directorio en el directorio recien leído
		avd.CreateDir(from, to[1:], pAvd)
		//Actualiza el directorio
		avd.WriteAvd(pAvd)
	} else if tipe == 2 {
		//El directorio no existe
		//Busca espacio en el bimap
		bm := Getbitmap(BitmapAvd)
		pAvd, _ := bm.FindSpaces(1)
		//Añade el puntero al avd
		a.addPointer(pAvd)
		//Escribe el nuevo directorio
		var avd Avd
		avd.NewAvd(to[0], "664", logUser.uid, logUser.gid)
		avd.WriteAvd(pAvd)
		writeBM(BitmapAvd, pAvd)
	}
}

//addPointer busca una posición libre para añadir el puntero n.
//Si ya no hay espacios verifica el apuntador indirecto
func (a *Avd) addPointer(n int32) {
	for _, pAvd := range a.ApArraySubdirectorios {
		if pAvd == 0 {
			pAvd = n
			return
		}
	}
	//No se escribió ninguno
	//Verifica el apuntado indirecto
	if a.ApArbolVirtualDirectorio == -1 {
		//No ha sido creado
		///Copia los atributos
		var avd Avd
		avd.copyAvd(*a)
		//Busca espacio en el bitmap
		bm := Getbitmap(BitmapAvd)
		pAvd, _ := bm.FindSpaces(1)
		//Añade el puntero n
		avd.addPointer(n)
		//Escribe el nuevo avd
		avd.WriteAvd(pAvd)
		//Escribe el bitmap
		writeBM(BitmapAvd, pAvd)
		a.ApArbolVirtualDirectorio = pAvd
	} else {
		//ya ha sido creado
		var avd Avd
		avd.ReadAvd(a.ApArbolVirtualDirectorio)
		avd.addPointer(n)
		avd.WriteAvd(a.ApArbolVirtualDirectorio)
	}
}

//copyAvd copia los atributos del avd
func (a *Avd) copyAvd(ref Avd) {
	copy(a.FechaCreacion[:], ref.FechaCreacion[:])
	copy(a.NombreDirectorio[:], ref.NombreDirectorio[:])
	copy(a.Auth[:], ref.Auth[:])
	a.ApArbolVirtualDirectorio = -1
	a.ApDetalleDirectorio = -1
	a.Proper = ref.Proper
	a.Gid = ref.Gid
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

//CreateRep recorre todos los registros de avd para generar un
//texto con nodos de grapvhiz
func (a *Avd) CreateRep(n int32) string {
	var strd strings.Builder
	//Recuper la representación en graphviz del avd actual
	strd.WriteString(a.getHTML(n))
	//Recorre el array de subdirectorios
	for _, pAvd := range a.ApArraySubdirectorios {
		if pAvd > 0 {
			var newAvd Avd
			newAvd.ReadAvd(pAvd)
			strd.WriteString(newAvd.CreateRep(pAvd))
		}
	}
	//Recupera el texto del aputnador indirecto
	if a.ApArbolVirtualDirectorio != -1 {
		var newAvd Avd
		newAvd.ReadAvd(a.ApArbolVirtualDirectorio)
		strd.WriteString(newAvd.CreateRep(a.ApArbolVirtualDirectorio))
	}
	return strd.String()
}

//getHTML recupera una tabla html que representa al avd en graphviz
func (a *Avd) getHTML(n int32) string {
	var strd strings.Builder
	var edge strings.Builder
	strd.WriteString("avd_" + fmt.Sprint(n) + " [\n")
	strd.WriteString("	label=<\n")
	strd.WriteString("		<table border='0' cellspacing='0' cellborder='1'>\n")
	//Recupera un slice del nombre de la partición hasta encontrar un caracter nulo
	idxEnd := bytes.IndexByte(a.NombreDirectorio[:], 0)
	if idxEnd == -1 {
		//Si no hay caracter nulo se tomará todo el array
		idxEnd = 20
	}
	tempName := string(a.NombreDirectorio[:idxEnd])
	strd.WriteString("		<tr>\n")
	strd.WriteString("		<td bgcolor='#2980b9' colspan='2'>" + tempName + "</td>\n")
	strd.WriteString("		</tr>\n")
	for i, pAvd := range a.ApArraySubdirectorios {
		if pAvd > 0 {
			strd.WriteString("		<tr>\n")
			strd.WriteString("		<td port='dp_" + fmt.Sprint(pAvd) + "' bgcolor='#2980b9' >dp" + fmt.Sprint(i+1) + "</td><td>" + fmt.Sprint(pAvd) + "</td>\n")
			edge.WriteString("avd_" + fmt.Sprint(n) + ":dp_" + fmt.Sprint(pAvd) + "-> avd_" + fmt.Sprint(pAvd) + "\n")
			strd.WriteString("		</tr>\n")
		}
	}
	if a.ApArbolVirtualDirectorio != -1 {
		strd.WriteString("		<tr>\n")
		strd.WriteString("		<td port='dp_" + fmt.Sprint(a.ApArbolVirtualDirectorio) + "' bgcolor='#2980b9'>ip</td><td>" + fmt.Sprint(a.ApArbolVirtualDirectorio) + "</td>\n")
		edge.WriteString("avd_" + fmt.Sprint(n) + ":ip_" + fmt.Sprint(a.ApArbolVirtualDirectorio) + " -> avd_" + fmt.Sprint(a.ApArbolVirtualDirectorio) + "\n")
		strd.WriteString("		</tr>\n")
	}
	strd.WriteString("		</table>\n")
	strd.WriteString("	>]\n\n")
	strd.WriteString(edge.String())
	return strd.String()
}
