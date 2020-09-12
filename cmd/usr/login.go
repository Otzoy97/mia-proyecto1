package usr

import (
	"mia-proyecto1/disk"
	"mia-proyecto1/lwh"
	"strings"

	"github.com/fatih/color"
)

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
	//Verifica que la partición esté montada
	path, name := disk.FindImg(l.id)
	//Monta el disco
	lwh.MountVDisk(path, name)
	//Busca el archivo user.txt
	var avd lwh.Avd
	pointer, typ := avd.Find("/user.txt")
	if typ == 2 || typ == 0 {
		//Es un error, no se pudo recuperar el puntero
		color.New(color.FgHiYellow).Printf("Login: no se pudo leer user.txt (%v)\n", l.Row)
		color.New(color.FgHiRed, color.Bold).Println("Login fracasó")
		return
	}
	//Si se pudo recuperar el puntero de user.txt
	//Recupera el contenido de user.txt
	var userTxt lwh.Inodo
	if userTxt.ReadInodo(pointer) {
		//Lee el contenido del inodo
		text := userTxt.GetCont()
		//Realiza un split con '\n'
		txtSpli := strings.Split(text, "\n")
		//Busca el uid y el gid
		uid, group := findUser(txtSpli, l.usr, l.pwd)
		if uid == 0 {
			color.New(color.FgHiYellow).Printf("Login: el usuario '%v' no existe (%v)\n", l.usr, l.Row)
			color.New(color.FgHiRed, color.Bold).Println("Login fracasó")
			return
		}
		//Buscca el gid
		gid := findGroup(txtSpli, group)
		if gid == 0 {
			color.New(color.FgHiYellow).Printf("Login: el grupo '%v' del usuario '%v' no existe (%v)\n", group, l.usr, l.Row)
			color.New(color.FgHiRed, color.Bold).Println("Login fracasó")
			return
		}
		//Almacena el usuario
		lwh.Login(uid, gid)
	}
}
