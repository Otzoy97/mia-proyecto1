package usr

import (
	"mia-proyecto1/lwh"

	"github.com/fatih/color"
)

//Logout ...
type Logout struct {
	Row int
}

//AddOp ...
func (l *Logout) AddOp(s string, v interface{}) {
}

//Validate ...
func (l *Logout) Validate() bool {
	return true
}

//Run ...
func (l *Logout) Run() {
	if lwh.Logout() {
		color.New(color.FgHiGreen, color.Bold).Println("Logout: sesión finalizada")
		return
	}
	color.New(color.FgHiYellow).Printf("Logout: no hay ninguna sesión activa (%v)\n", l.Row)
	color.New(color.FgHiRed, color.Bold).Println("Logout fracasó", "")
}
