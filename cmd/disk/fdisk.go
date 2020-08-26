package cmdisk

import (
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"os"
	"sort"
	"unsafe"

	"github.com/fatih/color"
)

//Fdisk ...
type Fdisk struct {
	unit, size    uint32
	Row           int
	path, name    string
	typ, fit, del byte
	Oplst         map[string]interface{}
	exec          byte
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
			m.size = uint32(m.Oplst["size"].(int))
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
					color.New(color.FgHiYellow).Printf("Fdisk: unit debe ser 'b', 'k' o 'm' (%v)\n", m.Row)
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
		} else {
			m.path = m.Oplst["path"].(string)
		}
		if !bName {
			color.New(color.FgHiYellow).Printf("Fdisk: name no se encontró (%v)\n", m.Row)
			f = false
		} else {
			m.name = m.Oplst["name"].(string)
		}
	}
	if !f {
		color.New(color.FgRed, color.Bold).Println("Fdisk no se puede ejecutar")
		return false
	}
	return true
}

//Run crea una nueva partición
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
		if m.putPartition(&mbr) && mbr.WriteMbr(file) {
			color.New(color.FgHiGreen, color.Bold).Printf("Fdisk: se creó la partición '%v' en el disco '%v' (%v)\n", m.name, m.path, m.Row)
		} else {
			color.New(color.FgHiYellow).Printf("Fdisk: no se pudo crear la partición '%v' en el disco '%v' (%v)\n", m.name, m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Fdisk fracasó")
		}
	}
}

//PutPartition coloca la particion con nombre partitionName de tamaño partitionSize
//de tipo partitionType y con fit partitionFit. True si se logra colocar la partiticion, false si no
func (m *Fdisk) putPartition(mbr *disk.Mbr) bool {
	if m.typ == 'l' {
		return false
	}
	return m.primaryPartition(mbr)
}

//primaryPartition coloca una partición primaria o extendida con los atributos especificados
//Para colocar la partición se utiliza FirstFit
func (m *Fdisk) primaryPartition(mbr *disk.Mbr) bool {
	//Recupera un array de particiones y si ya hay
	//una partición extendida
	arrPar, flag := cmd.CreateArrPart(mbr)
	if cmd.CheckNames(&arrPar, m.name) {
		color.New(color.FgHiYellow, color.Bold).Printf("Fdisk: nombre '%v' duplicado en '%v' (%v)\n", m.name, m.path, m.Row)
		return false
	}
	//Si la partición a crear es extendida y flag es true, se rechaza la acción
	if flag && m.typ == 'e' {
		color.New(color.FgHiYellow).Printf("Fdisk: no es posible colocar más de 1 partición extendida (%v)\n", m.Row)
		return false
	}
	//Si el arreglo es de 4 se rechaza la acción
	if len(arrPar) >= 4 {
		color.New(color.FgHiYellow).Printf("Fdisk: no hay espacio disponible para más particiones (%v)\n", m.Row)
		return false
	}
	//Ordena los elementos utilizando el atributo PartStart
	sort.Sort(disk.ByPartStart(arrPar))
	var mPar map[uint32]uint32 = map[uint32]uint32{}
	//Recorre todo el arreglo menos el último elemento
	for i := 0; i < len(arrPar)-1; i++ {
		par := arrPar[i]
		mPar[par.PartStart+par.PartSize] = arrPar[i+1].PartStart - (par.PartStart + par.PartSize) - 1
	}
	if len(arrPar) > 0 {
		//Calcula el byte de inicio y el tamaño disponible para el ultimo elemento
		par := arrPar[len(arrPar)-1]
		mPar[par.PartStart+par.PartSize] = mbr.MbrTamanio - (par.PartStart + par.PartSize) - 1
	} else {
		//Si es la primera partición que se coloca
		mPar[uint32(unsafe.Sizeof(mbr))] = mbr.MbrTamanio - uint32(unsafe.Sizeof(mbr))
	}
	//Coloca la partición en la primera posición que quepa
	for startByte, freeSpace := range mPar {
		if freeSpace >= m.size {
			if mbr.MbrPartition1.PartStatus == '0' {
				copy(mbr.MbrPartition1.PartName[:], m.name)
				mbr.MbrPartition1.PartFit = m.fit
				mbr.MbrPartition1.PartSize = m.size
				mbr.MbrPartition1.PartStart = startByte
				mbr.MbrPartition1.PartStatus = '1'
				mbr.MbrPartition1.PartType = m.typ
			} else if mbr.MbrPartition2.PartStatus == '0' {
				copy(mbr.MbrPartition2.PartName[:], m.name)
				mbr.MbrPartition2.PartFit = m.fit
				mbr.MbrPartition2.PartSize = m.size
				mbr.MbrPartition2.PartStart = startByte
				mbr.MbrPartition2.PartStatus = '1'
				mbr.MbrPartition2.PartType = m.typ
			} else if mbr.MbrPartition3.PartStatus == '0' {
				copy(mbr.MbrPartition3.PartName[:], m.name)
				mbr.MbrPartition3.PartFit = m.fit
				mbr.MbrPartition3.PartSize = m.size
				mbr.MbrPartition3.PartStart = startByte
				mbr.MbrPartition3.PartStatus = '1'
				mbr.MbrPartition3.PartType = m.typ
			} else if mbr.MbrPartition4.PartStatus == '0' {
				copy(mbr.MbrPartition4.PartName[:], m.name)
				mbr.MbrPartition4.PartFit = m.fit
				mbr.MbrPartition4.PartSize = m.size
				mbr.MbrPartition4.PartStart = startByte
				mbr.MbrPartition4.PartStatus = '1'
				mbr.MbrPartition4.PartType = m.typ
			}
			goto Eureka
		}
	}
	color.New(color.FgHiYellow).Printf("Fdisk: espacio insuficiente (%v)\n", m.Row)
	return false
Eureka:
	return true
}
