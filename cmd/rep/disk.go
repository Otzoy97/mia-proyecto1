package rep

import (
	"bytes"
	"fmt"
	"mia-proyecto1/disk"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

//CreateDisk ...
func (m *Rep) CreateDisk(mbr *disk.Mbr) []byte {
	var strD strings.Builder
	parArr, _ := mbr.CreateArrPart()
	sort.Sort(disk.ByPartStart(parArr))
	strD.WriteString("digraph G{\n")
	strD.WriteString("graph[pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"]\n")
	strD.WriteString("node [shape=plain fontname=\"Arial\"]\n")
	strD.WriteString("rankdir=LR\n")
	strD.WriteString("disk [label=<\n")
	strD.WriteString("<table border='2' cellborder='2' color='black' cellspacing='3'>\n")
	strD.WriteString("<tr>")
	strD.WriteString("<td bgcolor=\"#557571\">MBR<br/>" + strconv.Itoa(int(unsafe.Sizeof(*mbr))) + " bytes</td>\n")
	prev := uint32(unsafe.Sizeof(*mbr))
	for _, par := range parArr {
		libre := par.PartStart - prev
		if libre > 0 {
			strD.WriteString("<td bgcolor=\"#f4f4f4\">Free space<br/>" + strconv.Itoa(int(libre)) + " bytes</td>\n")
		}
		//Recupera un slice del nombre de la partición hasta encontrar un caracter nulo
		idxEnd := bytes.IndexByte(par.PartName[:], 0)
		if idxEnd == -1 {
			//Si no hay caracter nulo se tomará todo el array
			idxEnd = 16
		}
		tempName := par.PartName[:idxEnd]
		if par.PartType == 'p' {
			strD.WriteString("<td bgcolor=\"#f7d1ba\">Primary<br/>")
		} else if par.PartType == 'e' {
			strD.WriteString("<td bgcolor=\"#d49a89\">Extended<br/>")
		}
		strD.WriteString("<b>" + string(tempName) + "</b><br/>" + fmt.Sprint(par.PartSize) + " bytes</td>")
		prev = par.PartStart + par.PartSize
	}
	//Coloca el úlitmo espacio disponible
	if mbr.MbrTamanio-prev >= 0 {
		strD.WriteString("<td bgcolor=\"#f4f4f4\">Free space<br/>" + fmt.Sprint(mbr.MbrTamanio-prev) + " bytes</td>\n")
	}
	strD.WriteString("</tr>")
	strD.WriteString("</table>>]}\n")
	return []byte(strD.String())
}
