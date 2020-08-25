package disk

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

//Imgdsk ...
type Imgdsk struct {
	path  string
	autID int
	parts map[int]string
}

//letters used
var idxLetter map[byte]bool = map[byte]bool{
	'a': false,
	'b': false,
	'c': false,
	'd': false,
	'e': false,
	'f': false,
	'g': false,
	'h': false,
	'i': false,
	'j': false,
	'k': false,
	'l': false,
	'm': false,
	'n': false,
	'o': false,
	'p': false,
	'q': false,
	'r': false,
	's': false,
	't': false,
	'u': false,
	'v': false,
	'w': false,
	'x': false,
	'y': false,
	'z': false}

//mapa de imgdsk
var imglst map[byte]Imgdsk = map[byte]Imgdsk{}

//AddImg ...
func AddImg(path, name string) (bool, string) {
	//Verifica que la partición no esté ya montada
	for key, value := range imglst {
		//Si el path es el mismo
		if value.path == path {
			//Recorre los nombres
			for _, dskName := range value.parts {
				//Si el nombre coincide
				if dskName == name {
					//Ya está montada
					color.New(color.FgHiYellow).Printf("Mount: la particion '%v' del disco '%v' ya está montada\n", path, name)
					return false, ""
				}
			}
			//El ciclo terminó, no existe el nombre. Se añade el nuevo nombre
			value.autID++
			value.parts[value.autID] = name
			return true, "vd" + string(key) + string(value.autID)
		}
	}
	//No existe el nombre
	//Busca una letra disponible
	var letter byte
	for k, v := range idxLetter {
		if !v {
			letter = k
			goto NoReturn
		}
	}
	//Ya no se pueden montar más particiones
	color.New(color.FgHiYellow).Println("Mount: no se pueden montar más particiones")
	return false, ""
NoReturn:
	//Actualiza idxLetter
	idxLetter[letter] = true
	//Coloca el path y el name en el diccionari imglst
	imglst[letter] = Imgdsk{path: path,
		autID: 1,
		parts: map[int]string{1: name}}
	return true, "vd" + string(letter) + string(1)
}

//RmImg busca y elimina el registro de la partición
//Unmount utiliza esta función
func RmImg(id string) bool {
	//Verifica que el id tenga sentido
	idLow := strings.ToLower(id)
	if match, _ := regexp.Match(`vd[a-z][0-9]+`, []byte(idLow)); !match {
		color.New(color.FgHiYellow).Printf("Unmount: el id dado '%v' no es válido\n", id)
		return false
	}
	//Recupera la letra y el numero
	idLetter := []byte(idLow[2:3])
	idIdx, _ := strconv.Atoi(idLow[3:len(idLow)])
	//Busca en los diccionarios con las respectivas id
	if path, pExist := imglst[idLetter[0]]; pExist {
		//Busca el número
		if _, nExist := path.parts[idIdx]; nExist {
			//Elimina el registro
			delete(path.parts, idIdx)
			//Si ya no hay más regitros en parts, elimina el struct
			if len(path.parts) == 0 {
				delete(imglst, idLetter[0])
				//Actualiza las letras usadas
				idxLetter[idLetter[0]] = false
			}
			color.New(color.FgHiGreen, color.Bold).Printf("Unmount: se ha desmontado '%v'", id)
			return true
		}
	}
	color.New(color.FgHiYellow).Printf("Unmount: el id '%v' no existe\n", id)
	return false
}

//ListImg lista todas las particiones montadas
func ListImg() {
	headerFmt := color.New(color.FgHiWhite, color.Bold, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgBlue).SprintfFunc()
	tbl := table.New("ID", "PATH", "NAME")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	//Recorre el mapa y lista el nombre
	for letter, part := range imglst {
		for idx, name := range part.parts {
			str := "vd" + string(letter) + string(idx)
			tbl.AddRow(str, part.path, name)
		}
	}
	tbl.Print()
}

//FindImg recupera el path y el nombre de la partición
//Regresa path y name
func FindImg(id string) (string, string) {
	//Verifica que el id tenga sentido
	idLow := strings.ToLower(id)
	if match, _ := regexp.Match(`vd[a-z][0-9]+`, []byte(idLow)); !match {
		color.New(color.FgHiYellow).Printf("	El id dado '%v' no es válido\n", id)
		return "", ""
	}
	//Recupera la letra y el numero
	idLetter := []byte(idLow[2:3])
	idIdx, _ := strconv.Atoi(idLow[3:len(idLow)])
	//Busca en los diccionarios con las respectivas id
	if path, pExist := imglst[idLetter[0]]; pExist {
		//Busca el número
		if name, nExist := path.parts[idIdx]; nExist {
			//Devuelve el registro
			return path.path, name
		}
	}
	color.New(color.FgHiYellow).Printf("	El id '%v' no existe\n'", id)
	return "", ""
}
