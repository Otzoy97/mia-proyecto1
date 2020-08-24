package cmdisk

import (
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"os"

	"github.com/fatih/color"
)

//Fdisk ...
type Fdisk struct {
	unit, size, Row uint32
	path, name      string
	typ, fit, del   byte
	Oplst           map[string]interface{}
	exec            byte
}

//AddOp ...
func (m *Fdisk) AddOp(s string, v interface{}) {
	m.Oplst[s] = v
}

//Validate ...
func (m *Fdisk) Validate() bool {
	f := true
	bSize := cmd.ValidateOptions(&m.Oplst, "size")
	bDel := cmd.ValidateOptions(&m.Oplst, "delete")
	if !bSize && !bDel {
		color.New(color.FgHiYellow).Printf("Fdisk: size y delete no se encontró (%v)\n", m.Row)
		f = false
	} else {
		bPath := cmd.ValidateOptions(&m.Oplst, "path")
		bName := cmd.ValidateOptions(&m.Oplst, "name")
		if bSize {
			bUnit := cmd.ValidateOptions(&m.Oplst, "unit")
			bType := cmd.ValidateOptions(&m.Oplst, "type")
			bFit := cmd.ValidateOptions(&m.Oplst, "fit")
			m.unit = 1000
			m.fit = 'w'
			m.typ = 'p'
			m.size = m.Oplst["size"].(uint32)
			m.exec = 's'
			//Intenta recuperar elvalor de unit
			if bUnit {
				switch m.Oplst["unit"].(string) {
				case "k":
					m.unit = 1000
				case "m":
					m.unit = 1000000
				case "b":
					m.unit = 1
				default:
					color.New(color.FgHiYellow).Println("Fdisk: unit debe ser 'b', 'k' o 'm' (%v)\n", m.Row)
					f = false
				}
			}
			//Intenta recuperar el valor de fit
			if bFit {
				switch m.Oplst["fit"].(string) {
				case "ff":
					m.fit = 'f'
				case "wf":
					m.unit = 'w'
				case "bf":
					m.unit = 'b'
				default:
					color.New(color.FgHiYellow).Printf("Fdisk: fit debe ser 'ff', 'wf' o 'bf' (%v)\n", m.Row)
					f = false
				}
			}
			//Intenta recuperar el valor de type
			if bType {
				switch m.Oplst["type"].(string) {
				case "p":
					m.typ = 'p'
				case "e":
					m.typ = 'e'
				case "l":
					m.typ = 'l'
				default:
					color.New(color.FgHiYellow).Printf("Fdisk: unit debe ser 'p', 'e' o 'l' (%v)\n", m.Row)
					f = false
				}
			}
		} else if bDel {
			switch m.Oplst["delete"].(string) {
			case "fast":
				m.del = 'a'
			case "full":
				m.del = 'u'
			default:
				color.New(color.FgHiYellow).Printf("Fdisk: delete debe ser 'fast' o 'full (%v)\n", m.Row)
				f = false
			}
			m.exec = 'd'
		}
		if !bPath {
			color.New(color.FgHiYellow).Printf("Fdisk: path no se encontró (%v)\n", m.Row)
			f = false
		}
		if !bName {
			color.New(color.FgHiYellow).Printf("Fdisk: name no se encontró (%v)\n", m.Row)
			f = false
		}
	}
	if !f {
		color.New(color.FgRed, color.Bold).Printf("Fdisk no se puede ejecutar (%v)\n", m.Row)
		return false
	}
	return true
}

//Run ...
func (m *Fdisk) Run() {
	//Crear una partición
	//Verifica que el disco exista
	if _, err := os.Stat(m.path); err != nil {
		color.New(color.FgHiYellow).Printf("Fdisk: el disco '%v' no existe (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Printf("Fdisk fracasó (%v)\n", m.Row)
		return
	}
	//Abre el disco
	file, err := os.OpenFile(m.path, os.O_RDWR, 0777)
	if err != nil {
		//No se pudo abrir el disco
		color.New(color.FgHiYellow).Printf("Fdisk: no se pudo recuperar el disco '%v' (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Printf("Fdisk fracasó (%v)\n", m.Row)
		return
	}
	//Recupera el MBR
	mbr := disk.Mbr{}
	if !mbr.ReadMbr(file) {
		//No se pudo recuperar el mbr
		color.New(color.FgHiYellow).Printf("Fdisk: no se pudo recuperar el MBR '%v' (%v)\n", m.path, m.Row)
		color.New(color.FgHiRed, color.Bold).Printf("Fdisk fracasó (%v)\n", m.Row)
		return
	}
	switch m.exec {
	case 'd':
		//Eliminar una partición
		switch m.del {
		case 'a':
			//Eliminar mbr
		case 'u':
			//Eliminar mbr y limpiar espacio de partición
		}
	case 's':

	}
}
