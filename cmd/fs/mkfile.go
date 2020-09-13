package fs

import "strconv"

//Mkfile ...
type Mkfile struct {
	Row            int
	id, path, cont string
	recursive      bool
	size           int32
}

//AddOp ...
func (m *Mkfile) AddOp(s string, v interface{}) {
	if s == "id" {
		m.id = v.(string)
	} else if s == "path" {
		m.path = v.(string)
	} else if s == "p" {
		m.recursive = true
	} else if s == "size" {
		conv, _ := strconv.Atoi(v.(string))
		m.size = int32(conv)
	} else if s == "cont" {
		m.cont = v.(string)
	}
}

//Validate ...
func (m *Mkfile) Validate() bool {
	flag := true
	if !flag {
		return false
	}
	return true
}

//Run crea un nuevo archivo
func (m *Mkfile) Run() {

}
