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
var idxLetter [26]byte = [26]byte{}

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
					color.New(color.FgHiYellow).Printf("Mount: la particion '%v' del disco '%v' ya está montada\n", name, path)
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
	var letter byte = 'a'
	for _, v := range idxLetter {
		if v != letter {
			break
		}
		letter++
	}
	if letter > 'z' {
		//Ya no se pueden montar más particiones
		color.New(color.FgHiYellow).Println("Mount: no se pueden montar más particiones")
		return false, ""
	}
	//Actualiza idxLetter
	idxLetter[letter-97] = letter
	//Coloca el path y el name en el diccionari imglst
	imglst[letter] = Imgdsk{path: path,
		autID: 1,
		parts: map[int]string{1: name}}
	return true, "vd" + string(letter) + "1"
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
	idLetter := []byte(idLow[2:3])[0]
	idIdx, _ := strconv.Atoi(idLow[3:len(idLow)])
	//Busca en los diccionarios con las respectivas id
	if path, pExist := imglst[idLetter]; pExist {
		//Busca el número
		if _, nExist := path.parts[idIdx]; nExist {
			//Elimina el registro
			delete(path.parts, idIdx)
			//Si ya no hay más regitros en parts, elimina el struct
			if len(path.parts) == 0 {
				delete(imglst, idLetter)
				//Actualiza las letras usadas
				idxLetter[97-idLetter] = 0
			}
			color.New(color.FgHiGreen, color.Bold).Printf("Unmount: se ha desmontado '%v'\n", id)
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
			str := "vd" + string(letter) + strconv.Itoa(idx)
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

//CheckIfMounted recorre las particiones montadas y verifica con el path y el name si ya está montada
func CheckIfMounted(path, name string) bool {
	//Verifica si la partición está montada
	for _, value := range imglst {
		//Si el path es el mismo
		if value.path == path {
			//Recorre los nombres
			for _, dskName := range value.parts {
				//Si el nombre conincide
				if dskName == name {
					//Está montada
					return true
				}
			}
		}
	}
	return false
}
