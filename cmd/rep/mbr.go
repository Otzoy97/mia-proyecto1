package rep

import (
	"bytes"
	"fmt"
	"mia-proyecto1/disk"
	"strings"
	"time"
)

//CreateMBR crea el texto para escribir en el archivo dot
func (m *Rep) CreateMBR(mbr *disk.Mbr) []byte {
	var strD strings.Builder
	strD.WriteString("digraph G {\n")
	strD.WriteString("graph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"]\n")
	strD.WriteString("node [shape=plain]\n")
	strD.WriteString("rankdir = LR\n")
	strD.WriteString("mbr[\n")
	strD.WriteString("shape=plaintext \n")
	strD.WriteString("color=black \n")
	strD.WriteString("label=<\n")
	strD.WriteString("<table border='1' color='black' cellspacing='0' cellborder='1'>\n")
	strD.WriteString("<th>\n")
	strD.WriteString("<td>Nombre</td>\n")
	strD.WriteString("<td>Valor</td>\n")
	strD.WriteString("</th>\n")
	//----
	strD.WriteString("<tr>\n")
	strD.WriteString("<td> mbr_tamanio</td>\n")
	strD.WriteString("<td>" + fmt.Sprint(mbr.MbrTamanio) + "</td>\n")
	strD.WriteString("</tr>\n")
	//----
	strD.WriteString("<tr>\n")
	strD.WriteString("<td>mbr_fecha_creacion</td>\n")
	var tm time.Time
	tm.GobDecode(mbr.MbrFechaCreacion[:])
	strD.WriteString("<td>" + tm.Format("02 Jan 2006 03:04:05 PM") + "</td>\n")
	strD.WriteString("</tr>\n")
	//----
	strD.WriteString("<tr>\n")
	strD.WriteString("<td>mbr_disk_signature</td>\n")
	strD.WriteString("<td>" + fmt.Sprint(mbr.MbrDiskSignature) + "</td>\n")
	strD.WriteString("</tr>\n")
	pArr := []disk.Partition{}
	pArr = append(pArr, mbr.MbrPartition1)
	pArr = append(pArr, mbr.MbrPartition2)
	pArr = append(pArr, mbr.MbrPartition3)
	pArr = append(pArr, mbr.MbrPartition4)
	for i, par := range pArr {
		strD.WriteString("<tr>")
		strD.WriteString("<td>part_name_" + fmt.Sprint(i+1) + "</td>")
		if par.PartStatus == 1 {
			temName := par.PartName[:bytes.IndexByte(par.PartName[:], 0)]
			strD.WriteString("<td>" + string(temName) + "</td>")
		} else {
			strD.WriteString("<td></td>")
		}
		strD.WriteString("</tr>\n")
		strD.WriteString("<tr>")
		strD.WriteString("<td>part_status_" + fmt.Sprint(i+1) + "</td>")
		strD.WriteString("<td>" + fmt.Sprint(par.PartStatus) + "</td>")
		strD.WriteString("</tr>\n")
		strD.WriteString("<tr>")
		strD.WriteString("<td>part_type_" + fmt.Sprint(i+1) + "</td>")
		if par.PartStatus == 1 {
			strD.WriteString("<td>" + string(par.PartType) + "</td>")
		} else {
			strD.WriteString("<td></td>")
		}
		strD.WriteString("</tr>\n")
		strD.WriteString("<tr>")
		strD.WriteString("<td>part_fit_" + fmt.Sprint(i+1) + "</td>")
		if par.PartStatus == 1 {
			strD.WriteString("<td>" + string(par.PartFit) + "</td>")
		} else {
			strD.WriteString("<td></td>")
		}
		strD.WriteString("</tr>\n")
		strD.WriteString("<tr>")
		strD.WriteString("<td>part_start_" + fmt.Sprint(i+1) + "</td>")
		strD.WriteString("<td>" + fmt.Sprint(par.PartStart) + "</td>")
		strD.WriteString("</tr>\n")
		strD.WriteString("<tr>")
		strD.WriteString("<td>part_size_" + fmt.Sprint(i+1) + "</td>")
		strD.WriteString("<td>" + fmt.Sprint(par.PartSize) + "</td>")
		strD.WriteString("</tr>\n")
	}
	strD.WriteString("</table>\n")
	strD.WriteString(">]}\n")
	return []byte(strD.String())
}
