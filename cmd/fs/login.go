package fs

import "github.com/fatih/color"

//Login ...
type Login struct {
	Row          int
	usr, pwd, id string
}

//AddOp añade un parametro a utilizar
func (l *Login) AddOp(k string, v interface{}) {
	if k == "usr" {
		l.usr = v.(string)
	} else if k == "pwd" {
		l.pwd = v.(string)
	} else if k == "id" {
		l.id = v.(string)
	}
}

//Validate verifica que todas las opciones obligatorias existan
func (l *Login) Validate() bool {
	flag := true
	if l.usr == "" {
		color.New(color.FgHiYellow).Printf("Login: usr no se encontró (%v)\n", l.Row)
		flag = false
	}
	if l.pwd == "" {
		color.New(color.FgHiYellow).Printf("Login: pwd no se encontró (%v)\n", l.Row)
		flag = false
	}
	if l.id == "" {
		color.New(color.FgHiYellow).Printf("Login: id no se encontró (%v)\n", l.Row)
		flag = false
	}
	if !flag {
		color.New(color.FgHiRed, color.Bold).Println("Login no se puede ejecutar")
		return false
	}
	return true
}

//Run ejecuta el login
func (l *Login) Run() {
	//TODO: recuperar todo el contenido del user.txt
	//TODO: realizar un split con '\n', recorrer los valores y
	//realizar un split con ',',  recupera el usuario y contraseña
	//realizar match con ambos
}
