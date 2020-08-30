package rep

import (
	"bytes"
	"fmt"
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

//CreateDisk ...
func (m *Rep) CreateDisk(mbr *disk.Mbr) []byte {
	var strD strings.Builder
	parArr, _ := cmd.CreateArrPart(mbr)
	sort.Sort(disk.ByPartStart(parArr))
	strD.WriteString("digraph G{\n")
	strD.WriteString("graph[pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"]\n")
	strD.WriteString("node [shape=plain]\n")
	strD.WriteString("rankdir=LR\n")
	strD.WriteString("disk [label=<\n")
	strD.WriteString("<table border='1' cellborder='1' color='black' cellspacing='1'>\n")
	strD.WriteString("<tr>")
	strD.WriteString("<td>MBR<br/>" + strconv.Itoa(int(unsafe.Sizeof(*mbr))) + " bytes</td>\n")
	prev := uint32(unsafe.Sizeof(*mbr))
	for _, par := range parArr {
		libre := par.PartStart - prev
		if libre > 0 {
			strD.WriteString("<td>Free space<br/>" + strconv.Itoa(int(libre)) + " bytes</td>\n")
		}
		tempName := par.PartName[:bytes.IndexByte(par.PartName[:], 0)]
		if par.PartType == 'p' {
			strD.WriteString("<td>Primary<br/>")
		} else if par.PartType == 'e' {
			strD.WriteString("<td>Extended<br/>")
		}
		strD.WriteString("<b>" + string(tempName) + "</b><br/>" + fmt.Sprint(par.PartSize) + " bytes</td>")
		prev = par.PartStart + par.PartSize
	}
	//Coloca el úlitmo espacio disponible
	if mbr.MbrTamanio-prev >= 0 {
		strD.WriteString("<td>Free space<br/>" + fmt.Sprint(mbr.MbrTamanio-prev) + " bytes</td>\n")
	}
	strD.WriteString("</tr>")
	strD.WriteString("</table>>]}\n")
	return []byte(strD.String())
}
