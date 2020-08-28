package cmd

import (
	"bytes"
	"fmt"
	"mia-proyecto1/disk"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

//Rep ...
type Rep struct {
	Oplst                  map[string]interface{}
	path, id, ruta, nombre string
	exec                   bool
	Row                    int
}

//AddOp ...
func (m *Rep) AddOp(k string, v interface{}) {
	m.Oplst[k] = v
}

//Validate ...
func (m *Rep) Validate() bool {
	f := true
	m.exec = false
	bRuta := ValidateOptions(&m.Oplst, "ruta")
	if ValidateOptions(&m.Oplst, "nombre") {
		//Si el nombre es igual a ls o file, bRuta sería obligatorio
		if m.Oplst["nombre"].(string) == "ls" || m.Oplst["nombre"].(string) == "file" {
			if bRuta {
				m.ruta = m.Oplst["ruta"].(string)
				m.exec = true
			} else {
				color.New(color.FgHiYellow).Printf("Rep: ruta no se encontró (%v)\n", m.Row)
				f = false
			}
		}
	} else {
		color.New(color.FgHiYellow).Printf("Rep: nombre no se encontró (%v)\n", m.Row)
		f = false
	}
	if !ValidateOptions(&m.Oplst, "path") {
		color.New(color.FgHiYellow).Printf("Rep: path no se encontró (%v)\n", m.Row)
		f = false
	}
	if !ValidateOptions(&m.Oplst, "id") {
		color.New(color.FgHiYellow).Printf("Rep: id no se encontró (%v)\n", m.Row)
		f = false
	}
	if !f {
		return false
	}
	m.nombre = m.Oplst["nombre"].(string)
	m.path = m.Oplst["path"].(string)
	m.id = m.Oplst["id"].(string)
	return true
}

//Run ...
func (m *Rep) Run() {
	//Recupera el mbr
	mbr := disk.Mbr{}
	diskFile, err := os.OpenFile(m.path, os.O_RDWR, os.ModePerm)
	defer diskFile.Close()
	if err != nil {
		color.New(color.FgHiYellow).Printf("Rep: el disco '%v' no existe o no se puede abrir (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rep fracasó")
		return
	}
	if !mbr.ReadMbr(diskFile) {
		color.New(color.FgHiYellow).Printf("Rep: no se pudo recuperar el mbr del disco '%v' (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Println("Rep fracasó")
		return
	}
	//Crea el archivo qu contendra el texto para generar el grafo
	gFile, _ := os.Create("doFile.dot")
	switch m.nombre {
	case "mbr":
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
		b := []byte(strD.String())
		gFile.Write(b)
		gFile.Close()
		cmd := exec.Command("dot", "-Tpng", "doFile.dot", "-o", "mbr.png")
		cmd.Run()
		//os.Remove("doFile.dot")
	case "disk":

	}
}
