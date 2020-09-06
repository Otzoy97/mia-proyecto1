package disk

//Partition ...
type Partition struct {
	PartStatus, PartType, PartFit byte
	PartStart, PartSize           uint32
	PartName                      [16]byte
}

//ByPartStart ...
type ByPartStart []Partition

//ByPartStart sort.Interface basado en el campo PartStart
func (a ByPartStart) Len() int           { return len(a) }
func (a ByPartStart) Less(i, j int) bool { return a[i].PartStart < a[j].PartStart }
func (a ByPartStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//Find busca y devuelve una particion
func (a *ByPartStart) Find(name string) (Partition, bool) {
	//Recupera un array con todas las particiones activas
	byteName := [16]byte{}
	copy(byteName[:], name)
	//Recorre el arreglo de particiones activas
	for _, part := range *a {
		//Si el nombre es igual devuelve esa particion
		if part.PartName == byteName {
			return part, false
		}
	}
	return Partition{}, false
}

//CheckNames verifica si el nombre ya existe
func (a *ByPartStart) Check(name string) bool {
	//Verifica si el nombre ya existe
	for _, p := range *a {
		byteName := [16]byte{}
		copy(byteName[:], name)
		if byteName == p.PartName {
			//El nombre existe
			return true
		}
		if p.PartType == 'e' {
			//Si la partición es extendida verificará todas las particiones lógicas
		}
	}
	//No existe el nombre
	return false
}
