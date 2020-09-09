package rep

import (
	"bytes"
	"fmt"
	"mia-proyecto1/lwh"
	"strings"
	"time"
)

//CreateSB ...
func (r *Rep) CreateSB(sb *lwh.Superboot) []byte {
	var strD strings.Builder
	strD.WriteString("digraph G {\n")
	strD.WriteString("graph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"]\n")
	strD.WriteString("node [shape=plain fontname=\"Arial\"]\n")
	strD.WriteString("rankdir = LR\n")
	strD.WriteString("mbr[\n")
	strD.WriteString("shape=plaintext \n")
	strD.WriteString("color=black \n")
	strD.WriteString("label=<\n")
	strD.WriteString("<table border='1' color='black' cellspacing='0' cellborder='1'>\n")
	strD.WriteString("<th>\n")
	strD.WriteString("<td bgcolor=\"#4e89ae\"><b>Nombre</b></td>\n")
	strD.WriteString("<td bgcolor=\"#4e89ae\"><b>Valor</b></td>\n")
	strD.WriteString("</th>\n")
	//Nombre hd
	//Arma el nombre del disco
	idxEnd := bytes.IndexByte(sb.SbNombreHd[:], 0)
	if idxEnd == -1 {
		idxEnd = 16
	}
	strD.WriteString(addRow("sb_nombre_hd", string(sb.SbNombreHd[:idxEnd])))
	strD.WriteString(addRow("sb_arbol_virtual_count", sb.SbArbolvirtualCount))
	strD.WriteString(addRow("sb_detalle_directorio_count", sb.SbDetalleDirectorioCount))
	strD.WriteString(addRow("sb_inodos_count", sb.SbInodosCount))
	strD.WriteString(addRow("sb_bloques_count", sb.SbBloquesCount))
	strD.WriteString(addRow("sb_arbol_virtual_free", sb.SbArbolVirtualFree))
	strD.WriteString(addRow("sb_detalle_directorio_free", sb.SbDetalleDirectorioFree))
	strD.WriteString(addRow("sb_inodos_free", sb.SbInodosFree))
	strD.WriteString(addRow("sb_bloques_free", sb.SbBloquesFree))
	var tm time.Time
	tm.GobDecode(sb.SbDateCreacion[:])
	strD.WriteString(addRow("sb_date_creacion", tm.Format("02 Jan 2006 03:04:05 PM")))
	tm.GobDecode(sb.SbDateUltimoMontaje[:])
	strD.WriteString(addRow("sb_date_ultimo_montaje", tm.Format("02 Jan 2006 03:04:05 PM")))
	strD.WriteString(addRow("sb_montajes_count", sb.SbMontajesCount))
	strD.WriteString(addRow("sb_ap_bitmap_arbol_directorio", sb.SbApBitMapArbolDirectorio))
	strD.WriteString(addRow("sb_ap_arbol_directorio", sb.SbApArbolDirectorio))
	strD.WriteString(addRow("sb_ap_bitmap_detalle_directorio", sb.SbApBitmapDetalleDirectorio))
	strD.WriteString(addRow("sb_ap_detalle_directorio", sb.SbApDetalleDirectorio))
	strD.WriteString(addRow("sb_ap_bitmap_tabla_inodo", sb.SbApBitMapaTablaInodo))
	strD.WriteString(addRow("sb_ap_tabla_inodo", sb.SbApTablaInodo))
	strD.WriteString(addRow("sb_ap_bitmap_bloques", sb.SbApBitmapBloques))
	strD.WriteString(addRow("sb_ap_bloques", sb.SbApBloques))
	strD.WriteString(addRow("sb_ap_log", sb.SbApLog))
	strD.WriteString(addRow("sb_size_struct_arbol_directorio", sb.SbSizeStructArbolDirectorio))
	strD.WriteString(addRow("sb_size_struct_detalle_directorio", sb.SbSizeStructDetalleDirectorio))
	strD.WriteString(addRow("sb_size_struct_inodo", sb.SbSizeStructInodo))
	strD.WriteString(addRow("sb_size_struct_bloque", sb.SbSizeStructBloque))
	strD.WriteString(addRow("sb_first_free_bit_arbol_directorio", sb.SbFirstFreeBitArbolDirectorio))
	strD.WriteString(addRow("sb_first_free_bit_detalle_directorio", sb.SbFirstFreeBitDetalleDirectorio))
	strD.WriteString(addRow("sb_first_free_bit_tabla_inodo", sb.SbFirstFreeBitTablaInodo))
	strD.WriteString(addRow("sb_first_free_bit_bloques", sb.SbFirstFreeBitBloques))
	strD.WriteString(addRow("sb_magic_num", sb.SbMagicNum))
	strD.WriteString("</table>\n>]}\n")
	return []byte(strD.String())
}

func addRow(nombre string, valor interface{}) string {
	var strd strings.Builder
	strd.WriteString("<th>\n")
	strd.WriteString("<td>" + nombre + "</td>\n")
	strd.WriteString("<td>" + fmt.Sprint(valor) + "</td>\n")
	strd.WriteString("</th>\n")
	return strd.String()
}
