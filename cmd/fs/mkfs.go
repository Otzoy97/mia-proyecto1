package mkfs

//Mkfs ...
type Mkfs struct {
	id   string
	tipo byte
}

//AddOp ...
func (m *Mkfs) AddOp(s string, v interface{}) {
	if s == "id" {
		m.id = v.(string)
	} else if s == "tipo" {
		m.tipo = v.(byte)
	}
}

//Validate ....
func (m *Mkfs) Validate() bool {
	if m.tipo == 0 {
		m.tipo = 'u'
	}
	if m.id == "" {
		return false
	}
	return true
}

//Run ...
func (m *Mkfs) Run() {
	// Busca que la partición esté montada
	// Recupera el archivo del disco y la partición especificada
	// Solo trabaja sobre particiones primarias

}
