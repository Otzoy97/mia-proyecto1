package lwh

import (
	"bytes"
	"encoding/binary"
	"sort"

	"github.com/fatih/color"
)

//Bitmap ...
type Bitmap []byte

//BmType ...
type BmType byte

//BmType ...
const (
	BitmapAvd   BmType = 1
	BitmapDd           = 2
	BitmapInodo        = 3
	BitmapBd           = 4
)

//FindSpaces busca la primera posición libre en el bitmap que cumpla con el Fit f
//y que indique el inicio de espacio para albergar la cantidad 'spaces' de structs
func (b *Bitmap) FindSpaces(spaces int) (int32, bool) {
	var dict map[int]int = map[int]int{}
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
	switch vdPartition.PartFit {
	case 'w':
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
	case 'f':
		//Debe encontrar el primer valor que se mayor o igual a spaces
		for _, v := range keys {
			//Busca el primer valor que sea mayor o igual a spaces
			if dict[v] >= spaces {
				return int32(v), true
			}
		}
	case 'b':
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
func Getbitmap(op BmType) Bitmap {
	//Recupera el puntero del bitmap y el tamaño
	point, size := whichBM(op)
	if point == -1 {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el bitmap %v\n", virtualDisk.Name())
		return Bitmap{}
	}
	//Coloca el punteo de disco en posición
	virtualDisk.Seek(int64(point), 0)
	b := make(Bitmap, size)
	//Lee el bit map
	if _, err := virtualDisk.Read(b); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el bitmap %v\n", virtualDisk.Name())
		return Bitmap{}
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

//writeB escribe un byte en la posición p, del bitmap op
func writeBM(op BmType, p int32) bool {
	//Decide desde donde empezar a escribir
	bmpAp, _ := whichBM(op)
	//Coloca el apuntador de disco en posición
	offset := int64(bmpAp + p)
	virtualDisk.Seek(offset, 0)
	//Escribe el byte en la posición dada
	var b byte = 1
	bin := new(bytes.Buffer)
	binary.Write(bin, binary.BigEndian, b)
	if _, err := virtualDisk.Write(bin.Bytes()); err != nil {
		return false
	}
	return true
}

//Actualiza el número de bitmaps disponible en el superboot
func updateSBBM(op BmType) {
	switch op {
	case BitmapAvd:
		vdSuperBoot.SbArbolVirtualFree--
	case BitmapBd:
		vdSuperBoot.SbBloquesFree--
	case BitmapDd:
		vdSuperBoot.SbDetalleDirectorioFree--
	case BitmapInodo:
		vdSuperBoot.SbInodosFree--
	}
}
