package rep

import (
	"mia-proyecto1/lwh"
	"strings"
)

//CreateDir crea el texto para generar el reporte de directorio
func (r *Rep) CreateDir() []byte {
	var strD strings.Builder
	strD.WriteString("digraph G {\n")
	strD.WriteString("		graph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"]\n")
	strD.WriteString("		node [shape=plain fontname=\"Arial\"]\n")
	strD.WriteString("		rankdir = LR\n")
	//Lee el avd de la raiz
	var root lwh.Avd
	root.ReadAvd(0)
	//Recupera el texto para el reporte
	strD.WriteString(root.CreateRep(0, false))
	strD.WriteString("}\n")
	return []byte(strD.String())
}
