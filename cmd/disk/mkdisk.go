package cmdisk

//Mkdisk ...
type Mkdisk struct {
	size, Row int
	path      string
	name      string
	unit      string
	oplst     map[string]interface{}
}

//AddOp ...
func (m Mkdisk) AddOp(key string, value interface{}) {
	m.oplst[key] = value
}

//Validate ...
func (m Mkdisk) Validate() bool {
	return true
}

//Run ...
func (m Mkdisk) Run() {
	return
}
