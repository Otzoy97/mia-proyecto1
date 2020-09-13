package lwh

import (
	"sort"

	"github.com/fatih/color"
)

//Bitmap ...
type Bitmap []byte

//BmType ...
type BmType byte

//Fit ...
type Fit byte

//BmType ...
const (
	BitmapAvd   BmType = 1
	BitmapDd           = 2
	BitmapInodo        = 3
	BitmapBd           = 4
)

//Fit
const (
	WorstFit Fit = 1
	BestFit  Fit = 2
	FirstFit Fit = 3
)

//FindSpaces busca la primera posición libre en el bitmap que cumpla con el Fit f
//y que indique el inicio de espacio para albergar la cantidad 'spaces' de structs
func (b *Bitmap) FindSpaces(f Fit, spaces int) (int32, bool) {
	var dict map[int]int
	var cte int
	//Busca los espacios libres y desde ahí cuenta cuántos espacios libres hay
	for i, val := range *b {
		if val == 0 {
			dict[cte+1]++
		} else {
			cte = i
		}
	}
	//Arreglo para almacenar lasllaves
	keys := make([]int, len(dict))
	i := 0
	//Almacena todas las llaves de dict
	for k := range dict {
		keys[i] = k
		i++
	}
	//Ordena el arreglo keys
	sort.Ints(keys)
	switch f {
	case WorstFit:
		//Debe encontrar el espacio más grande
		var idxValue int
		//Inicializa un valor muy pequeño
		minValue := 0
		for _, v := range keys {
			//Busca el valor más grande que sea mayor o igual al valor de spaces
			if minValue < dict[v] && dict[v] >= spaces {
				minValue = dict[v]
				idxValue = v
			}
		}
		return int32(idxValue), true
	case FirstFit:
		//Debe encontrar el primer valor que se mayor o igual a spaces
		for _, v := range keys {
			//Busca el primer valor que sea mayor o igual a spaces
			if dict[v] >= spaces {
				return int32(v), true
			}
		}
	case BestFit:
		//Debe buscar el espacio más pequeño que sea mayor o igual al valor de spaces
		var idxValue int
		//Inicializa un valor muy grande
		maxValue := int(vdSuperBoot.SbBloquesCount)
		for _, v := range keys {
			//Busca el valor más pequeño que más se acerce al valor de spaces
			if maxValue > dict[v] && dict[v] >= spaces {
				if dict[v] == spaces {
					//Si encuentra un valor igual, retorna inmediatamente ese valor
					return int32(v), true
				} else if dict[v] > spaces {
					idxValue = v
					maxValue = dict[v]
				}
			}
		}
		return int32(idxValue), true
	}
	return -1, false
}

//Getbitmap recupera el stream de bytes que representan los espacios
//llenos y vacios del respectivo tipo de bitmap
func Getbitmap(op BmType) []byte {
	//Recupera el puntero del bitmap y el tamaño
	point, size := whichBM(op)
	if point == -1 {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el bitmap %v\n", virtualDisk.Name())
		return []byte{}
	}
	//Coloca el punteo de disco en posición
	virtualDisk.Seek(int64(point), 0)
	b := make([]byte, size)
	//Lee el bit map
	if _, err := virtualDisk.Read(b); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el bitmap %v\n", virtualDisk.Name())
		return []byte{}
	}
	return b
}

//whichBM decide el puntero a utilizar y la cantidad de bytes a leer
func whichBM(t BmType) (int32, int32) {
	switch t {
	case BitmapAvd:
		return vdSuperBoot.SbApBitMapArbolDirectorio, vdSuperBoot.SbArbolvirtualCount
	case BitmapBd:
		return vdSuperBoot.SbApBitmapBloques, vdSuperBoot.SbBloquesCount
	case BitmapDd:
		return vdSuperBoot.SbApBitmapDetalleDirectorio, vdSuperBoot.SbDetalleDirectorioCount
	case BitmapInodo:
		return vdSuperBoot.SbApBitMapaTablaInodo, vdSuperBoot.SbInodosCount
	}
	return -1, 0
}
