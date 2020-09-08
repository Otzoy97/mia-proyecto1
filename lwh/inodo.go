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
	Gid                   int32
}

//NewInodo crea un nuevo inodo
func (i *Inodo) NewInodo(pInodo, proper, gid int, auth string) {
	i.CountInodo = int32(pInodo)
	i.SizeArchivo = 0
	i.CountBloquesAsignados = 0
	i.ArrayBloques = [4]int32{-1, -1, -1, -1}
	i.ApIndirecto = -1
	i.IDProper = int32(proper)
	i.Gid = int32(gid)
	copy(i.Auth[:], auth)
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
