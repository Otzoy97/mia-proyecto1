package analyzer

import (
	"fmt"
	"mia-proyecto1/cmd"
	cmdisk "mia-proyecto1/cmd/disk"
	"reflect"
)

//Regresa una función que acepta un arreglo s que inserta
//las cadenas especificadas en t
var setState = func(t ...interface{}) func(s *[]interface{}) {
	return func(s *[]interface{}) {
		*s = append(*s, t...)

	}
}

//Define la tabla de análisis sintáctico
var parserTable map[string]map[string]func(*[]interface{}) = map[string]map[string]func(*[]interface{}){
	"S": map[string]func(*[]interface{}){
		"mkfs":     setState("Cmdlst"),
		"login":    setState("Cmdlst"),
		"logout":   setState("Cmdlst"),
		"mkgrp":    setState("Cmdlst"),
		"mkusr":    setState("Cmdlst"),
		"mkfile":   setState("Cmdlst"),
		"mkdir":    setState("Cmdlst"),
		"loss":     setState("Cmdlst"),
		"recovery": setState("Cmdlst"),
		"rep":      setState("Cmdlst"),
		"unmount":  setState("Cmdlst"),
		"mount":    setState("Cmdlst"),
		"fdisk":    setState("Cmdlst"),
		"mkdisk":   setState("Cmdlst"),
		"pause":    setState("Cmdlst"),
		"exec":     setState("Cmdlst"),
		"$":        setState()},
	"Cmdlst": map[string]func(*[]interface{}){
		"mkfs":     setState("Cmdlst1", "Cmd"),
		"login":    setState("Cmdlst1", "Cmd"),
		"logout":   setState("Cmdlst1", "Cmd"),
		"mkgrp":    setState("Cmdlst1", "Cmd"),
		"mkusr":    setState("Cmdlst1", "Cmd"),
		"mkfile":   setState("Cmdlst1", "Cmd"),
		"mkdir":    setState("Cmdlst1", "Cmd"),
		"loss":     setState("Cmdlst1", "Cmd"),
		"recovery": setState("Cmdlst1", "Cmd"),
		"rep":      setState("Cmdlst1", "Cmd"),
		"unmount":  setState("Cmdlst1", "Cmd"),
		"mount":    setState("Cmdlst1", "Cmd"),
		"fdisk":    setState("Cmdlst1", "Cmd"),
		"mkdisk":   setState("Cmdlst1", "Cmd"),
		"pause":    setState("Cmdlst1", "Cmd"),
		"exec":     setState("Cmdlst1", "Cmd")},
	"Cmdlst1": map[string]func(*[]interface{}){
		"mkfs":     setState("Cmdlst1", "Cmd"),
		"login":    setState("Cmdlst1", "Cmd"),
		"logout":   setState("Cmdlst1", "Cmd"),
		"mkgrp":    setState("Cmdlst1", "Cmd"),
		"mkusr":    setState("Cmdlst1", "Cmd"),
		"mkfile":   setState("Cmdlst1", "Cmd"),
		"mkdir":    setState("Cmdlst1", "Cmd"),
		"loss":     setState("Cmdlst1", "Cmd"),
		"recovery": setState("Cmdlst1", "Cmd"),
		"rep":      setState("Cmdlst1", "Cmd"),
		"unmount":  setState("Cmdlst1", "Cmd"),
		"mount":    setState("Cmdlst1", "Cmd"),
		"fdisk":    setState("Cmdlst1", "Cmd"),
		"mkdisk":   setState("Cmdlst1", "Cmd"),
		"pause":    setState("Cmdlst1", "Cmd"),
		"exec":     setState("Cmdlst1", "Cmd"),
		"$":        setState()},
	"Cmd": map[string]func(*[]interface{}){
		"mkfs":     setState("Mkfs"),
		"login":    setState("Login"),
		"logout":   setState("Logout"),
		"mkgrp":    setState("Mkgrp"),
		"mkusr":    setState("Mkusr"),
		"mkfile":   setState("Mkfile"),
		"mkdir":    setState("Mkdir"),
		"loss":     setState("Loss"),
		"recovery": setState("Recovery"),
		"rep":      setState("Rep"),
		"unmount":  setState("Unmount"),
		"mount":    setState("Mount"),
		"fdisk":    setState("Fdisk"),
		"mkdisk":   setState("Mkdisk"),
		"pause":    setState("Pause"),
		"exec":     setState("Exec")},
	"Mkfs": map[string]func(*[]interface{}){
		"mkfs": setState("Oplst", "mkfs")},
	"Login": map[string]func(*[]interface{}){
		"login": setState("Oplst", "login")},
	"Logout": map[string]func(*[]interface{}){
		"logout": setState("logout")},
	"Mkgrp": map[string]func(*[]interface{}){
		"mkgrp": setState("Oplst", "mkgrp")},
	"Mkusr": map[string]func(*[]interface{}){
		"mkusr": setState("Oplst", "mkusr")},
	"Mkfile": map[string]func(*[]interface{}){
		"mkfile": setState("Oplst", "mkfile")},
	"Mkdir": map[string]func(*[]interface{}){
		"mkdir": setState("Oplst", "mkdir")},
	"Loss": map[string]func(*[]interface{}){
		"loss": setState("Oplst", "loss")},
	"Recovery": map[string]func(*[]interface{}){
		"recovery": setState("Oplst", "recovery")},
	"Rep": map[string]func(*[]interface{}){
		"rep": setState("Oplst", "rep")},
	"Unmount": map[string]func(*[]interface{}){
		"unmount": setState("Oplst", "unmount")},
	"Mount": map[string]func(*[]interface{}){
		"mount": setState("Oplst", "mount")},
	"Fdisk": map[string]func(*[]interface{}){
		"fdisk": setState("Oplst", "fdisk")},
	"Mkdisk": map[string]func(*[]interface{}){
		"mkdisk": setState("Oplst", "mkdisk")},
	"Pause": map[string]func(*[]interface{}){
		"pause": setState("pause")},
	"Exec": map[string]func(*[]interface{}){
		"exec": setState("Oplst", "exec")},
	"Oplst": map[string]func(*[]interface{}){
		"tipo":   setState("Oplst1", "Op"),
		"grp":    setState("Oplst1", "Op"),
		"pwd":    setState("Oplst1", "Op"),
		"usr":    setState("Oplst1", "Op"),
		"cont":   setState("Oplst1", "Op"),
		"p":      setState("Oplst1", "Op"),
		"ruta":   setState("Oplst1", "Op"),
		"nombre": setState("Oplst1", "Op"),
		"id":     setState("Oplst1", "Op"),
		"add":    setState("Oplst1", "Op"),
		"delete": setState("Oplst1", "Op"),
		"fit":    setState("Oplst1", "Op"),
		"type":   setState("Oplst1", "Op"),
		"unit":   setState("Oplst1", "Op"),
		"name":   setState("Oplst1", "Op"),
		"size":   setState("Oplst1", "Op"),
		"path":   setState("Oplst1", "Op")},
	"Oplst1": map[string]func(*[]interface{}){
		"tipo":     setState("Oplst1", "Op"),
		"grp":      setState("Oplst1", "Op"),
		"pwd":      setState("Oplst1", "Op"),
		"usr":      setState("Oplst1", "Op"),
		"cont":     setState("Oplst1", "Op"),
		"p":        setState("Oplst1", "Op"),
		"ruta":     setState("Oplst1", "Op"),
		"nombre":   setState("Oplst1", "Op"),
		"id":       setState("Oplst1", "Op"),
		"add":      setState("Oplst1", "Op"),
		"delete":   setState("Oplst1", "Op"),
		"fit":      setState("Oplst1", "Op"),
		"type":     setState("Oplst1", "Op"),
		"unit":     setState("Oplst1", "Op"),
		"name":     setState("Oplst1", "Op"),
		"size":     setState("Oplst1", "Op"),
		"path":     setState("Oplst1", "Op"),
		"mkfs":     setState(),
		"login":    setState(),
		"logout":   setState(),
		"mkgrp":    setState(),
		"mkusr":    setState(),
		"mkfile":   setState(),
		"mkdir":    setState(),
		"loss":     setState(),
		"recovery": setState(),
		"rep":      setState(),
		"unmount":  setState(),
		"mount":    setState(),
		"fdisk":    setState(),
		"mkdisk":   setState(),
		"pause":    setState(),
		"exec":     setState(),
		"$":        setState()},
	"Op": map[string]func(*[]interface{}){
		"tipo":   setState("opDel", "asignacion", "tipo"),
		"grp":    setState("cadena", "asignacion", "grp"),
		"pwd":    setState("cadena", "asignacion", "pwd"),
		"usr":    setState("cadena", "asignacion", "usr"),
		"cont":   setState("cadena", "asignacion", "cont"),
		"p":      setState("p"),
		"ruta":   setState("cadena", "asignacion", "ruta"),
		"nombre": setState("opRep", "asignacion", "nombre"),
		"id":     setState("cadena", "asignacion", "id"),
		"add":    setState("numero", "asignacion", "add"),
		"delete": setState("opDel", "asignacion", "delete"),
		"fit":    setState("opFit", "asignacion", "fit"),
		"type":   setState("opType", "asignacion", "type"),
		"unit":   setState("opUnit", "asignacion", "unit"),
		"name":   setState("cadena", "asignacion", "name"),
		"size":   setState("numero", "asignacion", "size"),
		"path":   setState("cadena", "asignacion", "path")}}

var cmdlst []cmd.Command

//Parser utiliza los token generados por el analizador léxico
func Parser() {
	//Apuntador para el flujo de tokens
	pToken := 0
	//Crea el stack
	stack := []interface{}{"$", "S"}
	//Agrega un $ al final del stream de tokens
	tokQueue = append(tokQueue, &Token{lex: "$", row: tokQueue[len(tokQueue)-1].row, col: tokQueue[len(tokQueue)-1].col + 1, tokname: "$"})
	//Ultimo elemento del stack
	for stack[len(stack)-1] != "$" {
		x := stack[len(stack)-1]
		if len(tokQueue)-1 < pToken {
			fmt.Println("Fin de tokens inesperado")
			return
		}
		if isTerminal(x) {
			//Extrae un token
			refX := tokQueue[pToken]
			pToken++
			if reflect.TypeOf(x).Kind() == reflect.Func {
				x.(func(int))(pToken)
			} else if x.(string) != refX.tokname {
				fmt.Printf("Se esperaba %v, se encontró %v (%v, %v)\n", x, refX.tokname, refX.row, refX.col)
				//Ejecuta un modo pánico
				continue
			} else {
				//Extrae el terminal de la pila
				stack = stack[:len(stack)-1]
			}
		} else {
			if f := parserTable[x.(string)][tokQueue[pToken].tokname]; f != nil {
				stack = stack[:len(stack)-1]
				parserActions(x.(string))
				f(&stack)
			} else {
				fmt.Print("Se esperaba ")
				for k := range parserTable[x.(string)] {
					fmt.Print(k, " ")
				}
				fmt.Printf("(%v, %v)\n", tokQueue[pToken].row, tokQueue[pToken].col)
				pToken++
			}
		}
	}
	//fmt.Println("Análisis sintáctico exitoso")
}

//Recibe una cadena y regresa una función
//El argumento t de la función retornada
//es la posición del apuntador de token al consumir el token
func parserActions(s string) func(t int) {
	switch s {
	case "S":
	case "Cmd":
	case "Cmdlst":
	case "Cmdlst1":
	case "Oplst":
	case "Oplst1":
	case "Exec":
	case "Pause":
	case "Mkdisk":
		return func(t int) {
			cmdlst = append(cmdlst, cmdisk.Mkdisk{Row: tokQueue[t].row})
		}
	case "Fdisk":
	case "Mount":
	case "Unmount":
	case "Rep":
	case "Recovery":
	case "Loss":
	case "Mkdir":
	case "Mkfile":
	case "Mkusr":
	case "Mkgrp":
	case "Logout":
	case "Login":
	case "Mkfs":
	case "Op":
		return func(t int) {
			if tokQueue[t].tokname == "p" {
				cmdlst[len(cmdlst)-1].AddOp("p", true)
			} else {
				cmdlst[len(cmdlst)-1].AddOp(tokQueue[t].tokname, tokQueue[t-2].lex)
			}
		}
	}
	return func(t int) {

	}
}

//Verifica si la cadena representa un terminal
func isTerminal(s interface{}) bool {
	if reflect.TypeOf(s).Kind() == reflect.Func {
		return false
	}
	//Si el primer caracter es 'mayúsucula', entonces es terminal
	sx := []rune(s.(string))
	if sx[0] >= 'A' && sx[0] <= 'Z' {
		return false
	}
	return true
}
