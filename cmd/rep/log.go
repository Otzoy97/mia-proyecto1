package rep

import (
	"bytes"
	"fmt"
	"mia-proyecto1/lwh"
	"strings"
	"time"
)

//CreateLog crea el texto para el reporte de la bitacora
func (r *Rep) CreateLog(lg []lwh.Log) []byte {
	var strd strings.Builder
	var tm time.Time
	//Recorre el array de la bitacora
	for i, v := range lg {
		//Traduce la fecha
		tm.GobDecode(v.LogFecha[:])
		strd.WriteString(fmt.Sprint(i+1) + " - " + tm.Format("02 Jan 2006 03:04:05 PM") + "\n")
		strd.WriteString("Operación: " + v.Getop() + "\n")
		strd.WriteString("Tipo: " + v.Getipo() + "\n")
		//Recupera el nombre
		idxEnd := bytes.IndexByte(v.LogNombre[:], 0)
		if idxEnd == -1 {
			idxEnd = 256
		}
		strd.WriteString("Nombre: " + string(v.LogNombre[:idxEnd]) + "\n")
		//Recupera el contenido
		idxEnd = bytes.IndexByte(v.LogContenido[:], 0)
		if idxEnd == -1 {
			idxEnd = 256
		}
		strd.WriteString("Contenido: \"" + string(v.LogContenido[:idxEnd]) + "\"\n\n")
	}
	return []byte(strd.String())
}
