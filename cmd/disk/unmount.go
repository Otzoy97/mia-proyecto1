package cmdisk

import (
	"mia-proyecto1/disk"

	"github.com/fatih/color"
)

//Unmount ...
type Unmount struct {
	Row   int
	idlst []string
}

//AddOp agrega un nuevo id, si no es un id lo rechaza
func (m *Unmount) AddOp(k string, v interface{}) {
	if "id" == k {
		m.idlst = append(m.idlst, v.(string))
	}
}

//Validate ...
func (m *Unmount) Validate() bool {
	if len(m.idlst) > 0 {
		return true
	}
	color.New(color.FgHiYellow).Printf("Unmount no se puede ejecutar (%v)\n", m.Row)
	return false
}

//Run desmonta todas las particiones listadas
func (m *Unmount) Run() {
	for _, id := range m.idlst {
		disk.RmImg(id)
	}
}
