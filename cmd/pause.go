package cmd

import (
	"bufio"
	"fmt"
	"os"
)

//Pause ...
type Pause struct {
}

//AddOp ...
func (m *Pause) AddOp(s string, i interface{}) {
	return
}

//Validate ...
func (m *Pause) Validate() bool {
	return true
}

//Run ...
func (m *Pause) Run() {
	fmt.Print("Presiona 'Enter' para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return
}
