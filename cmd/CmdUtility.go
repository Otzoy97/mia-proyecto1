package cmd

import (
	"github.com/fatih/color"
)

//ValidateOptions verifica que el parametro s exista en el mapa m
//se asegura que el tipo de dato sea el correcto
func ValidateOptions(m *map[string]interface{}, s string) bool {
	//Verifica que la llave exista
	if _, mCheck := (*m)[s]; mCheck {
		return false
	}
	//Recupera el valor de la llave y el tipo de dato
	mVal := (*m)[s]
	//Verifica que el tipo y el valor de mVal sea el correcto
	switch s {
	case "tipo", "delete":
		if mVal != "fast" && mVal != "full" {
			color.New(color.FgHiYellow).Printf("El parámetro '%v' se ignoró. '%v' no se reconoció.\n", s, mVal)
			color.New(color.FgHiYellow).Printf("Parámetro '%v' debe tener como valor %v o %v.\n", s, "fast", "full")
			return false
		}
		return true
	case "grp", "pwd", "usr", "cont", "ruta", "id", "name", "path", "p":
		return true
	case "nombre":
		if mVal == "mbr" ||
			mVal == "disk" ||
			mVal == "sb" ||
			mVal == "bm_arbdir" ||
			mVal == "bm_detdir" ||
			mVal == "bm_inode" ||
			mVal == "bm_block" ||
			mVal == "bitacora" ||
			mVal == "directorio" ||
			mVal == "tree_file" ||
			mVal == "tree_directorio" ||
			mVal == "tree_complete" ||
			mVal == "ls" {
			return true
		}
		color.New(color.FgHiYellow).Printf("El parámetro 'nombre' se ignoró. '%v' no se reconoció.\n", mVal)
		color.New(color.FgHiYellow).Print("Parámetro 'nombre' debe tener como valor ")
		color.New(color.FgHiYellow).Printf("%v, %v, %v, %v", "mbr", "disk", "sb", "bm_arbdir")
		color.New(color.FgHiYellow).Printf("%v, %v, %v, %v", "bm_detdir", "bm_inode", "bm_block", "bitacora")
		color.New(color.FgHiYellow).Printf("%v, %v, %v o %v.\n", "tree_file", "tree_directorio", "tree_complete", "ls")
		return false
	case "add", "size":
		return true
	case "fit":
		if mVal == "bf" || mVal == "ff" || mVal == "wf" {
			return true
		}
		color.New(color.FgHiYellow).Printf("El parámetro 'fit' se ignoró. '%v' no se reconoció.\n", mVal)
		color.New(color.FgHiYellow).Printf("Parámetro 'fit' debe tener como valor %v, %v, %v\n", "bf", "ff", "wf")
		return false
	case "type":
		if mVal == "p" || mVal == "e" || mVal == "l" {
			return true
		}
		color.New(color.FgHiYellow).Printf("El parámetro 'type' se ignoró. '%v' no se reconoció.\n", mVal)
		color.New(color.FgHiYellow).Printf("Parámetro 'type' debe tener como valor %v, %v, %v\n", "bf", "ff", "wf")
		return false
	case "unit":
		if mVal == "m" || mVal == "k" || mVal == "b" {
			return true
		}
		color.New(color.FgHiYellow).Printf("El parámetro 'unit' se ignoró. '%v' no se reconoció.\n", mVal)
		color.New(color.FgHiYellow).Printf("Parámetro 'unit' debe tener como valor %v, %v, %v\n", "m", "k", "b")
		return false
	}
	return false
}
