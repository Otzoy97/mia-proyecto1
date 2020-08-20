package analyzer

import "fmt"

//Regresa una función que acepta un arreglo s que inserta
//las cadenas especificadas en t
var setState = func(t ...string) func(s *[]string) {
	return func(s *[]string) {
		*s = append(*s, t...)
	}
}

//Define la tabla de análisis sintáctico
var parserTable map[string]map[string]func(*[]string) = map[string]map[string]func(*[]string){
	"S": map[string]func(*[]string){
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
	"Cmdlst": map[string]func(*[]string){
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
	"Cmdlst1": map[string]func(*[]string){
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
	"Cmd": map[string]func(*[]string){
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
	"Mkfs": map[string]func(*[]string){
		"mkfs": setState("mkfs")},
	"Login": map[string]func(*[]string){
		"login": setState("login")},
	"Logout": map[string]func(*[]string){
		"logout": setState("logout")},
	"Mkgrp": map[string]func(*[]string){
		"mkgrp": setState("mkgrp")},
	"Mkusr": map[string]func(*[]string){
		"mkusr": setState("mkusr")},
	"Mkfile": map[string]func(*[]string){
		"mkfile": setState("mkfile")},
	"Mkdir": map[string]func(*[]string){
		"mkdir": setState("mkdir")},
	"Loss": map[string]func(*[]string){
		"loss": setState("loss")},
	"Recovery": map[string]func(*[]string){
		"recovery": setState("recovery")},
	"Rep": map[string]func(*[]string){
		"rep": setState("rep")},
	"Unmount": map[string]func(*[]string){
		"unmount": setState("unmount")},
	"Mount": map[string]func(*[]string){
		"mount": setState("mount")},
	"Fdisk": map[string]func(*[]string){
		"fdisk": setState("fdisk")},
	"Mkdisk": map[string]func(*[]string){
		"mkdisk": setState("mkdisk")},
	"Pause": map[string]func(*[]string){
		"pause": setState("pause")},
	"Exec": map[string]func(*[]string){
		"exec": setState("exec")},
	"Oplst": map[string]func(*[]string){
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
	"Oplst1": map[string]func(*[]string){
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
	"Op": map[string]func(*[]string){
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
		"size":   setState("cadena", "asignacion", "size"),
		"path":   setState("cadena", "asignacion", "path")}}

//Define el stack para el analizador LL(1)
var stack []string = []string{"$", "S"}

//Parser utiliza los token generados por el analizador léxico
func Parser() {
	//Agrega un $ al final del stream de tokens
	tokQueue = append(tokQueue, &Token{lex: "$", row: -1, col: -1, tokname: "$"})
	//Ultimo elemento del stack
	for stack[len(stack)-1] != "$" {
		x := stack[len(stack)-1]
		if len(tokQueue) < 1 {
			fmt.Println("Fin de tokens inesperado")
			return
		}
		if isTerminal(x) {
			//Extrae un token
			refX := tokQueue[0]
			tokQueue = tokQueue[1:]
			if x != refX.tokname {
				fmt.Printf("Se esperaba %v (%v, %v)\n", refX.tokname, refX.row, refX.col)
				//Ejecuta un modo pánico
				continue
			} else {
				//Extrae el no terminal de la pila
				stack = stack[:len(stack)-1]
			}
		} else {
			if f := parserTable[x][tokQueue[0].tokname]; f != nil {
				stack = stack[:len(stack)-1]
				f(&stack)
			} else {
				fmt.Print("Se esperaba ")
				for k := range parserTable[x] {
					fmt.Print(k, " ")
				}
				fmt.Printf("(%v, %v)\n", tokQueue[0].row, tokQueue[0].col)
				tokQueue = tokQueue[1:]
			}
		}
	}
	fmt.Println("Análisis sintáctico exitoso")
}

//Verifica si la cadena representa un terminal
func isTerminal(s string) bool {
	//Si el primer caracter es 'mayúsucula', entonces es terminal
	sx := []rune(s)
	if sx[0] >= 'A' && sx[0] <= 'Z' {
		return false
	}
	return true
}
