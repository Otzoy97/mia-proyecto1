package cmdisk

import "fmt"

//Mkdisk ...
type Mkdisk struct {
	size, Row int
	path      string
	name      string
	unit      string
	Oplst     map[string]interface{}
}

//AddOp ...
func (m Mkdisk) AddOp(key string, value interface{}) {
	m.Oplst[key] = value
}

//Validate ...
func (m Mkdisk) Validate() bool {
	return true
}

//Run ...
func (m Mkdisk) Run() {
	for k, v := range m.Oplst {
		fmt.Printf("%v -> %v\n", k, v)
	}
}
