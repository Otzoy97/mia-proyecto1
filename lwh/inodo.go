package lwh

//Inodo ...
type Inodo struct {
	CountInodo            int32
	SizeArchivo           int32
	CountBloquesAsignados int32
	ArrayBloques          [4]int32
	ApIndirecto           int32
	IDProper              int32
	Auth                  [3]byte
}

//getCont lee cada uno de los bloques y concatena el
//contenido que alojan
func (i *Inodo) getCont() string {
	//Viaja a la posición del inodo que especifica el
	//array de IArrayBloques y recupera el contenido binario
	// for _, offset := i.IArrayBloques{
	// 	//file.Seek()
	// }
	return ""
}
