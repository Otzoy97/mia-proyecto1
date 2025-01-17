package cmdisk

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"mia-proyecto1/cmd"
	"mia-proyecto1/disk"
	"os"
	"sort"
	"strings"
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
			m.unit = 1024
			m.fit = 'w'
			m.typ = 'p'
			m.exec = 's'
			//Intenta recuperar elvalor de unit
			if bUnit {
				switch m.Oplst["unit"].(string) {
				case "k":
					m.unit = 1024
				case "m":
					m.unit = 1024 * 1024
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
					m.fit = 'w'
				case "bf":
					m.fit = 'b'
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
			//guarda el tamaño
			if m.Oplst["size"].(int) < 0 {
				color.New(color.FgHiYellow).Printf("Fdisk: size debe ser mayor a 0 (%v)\n", m.Row)
				f = false
			} else {
				m.size = uint32(m.Oplst["size"].(int)) * m.unit
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
	defer file.Close()
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
		//Busca el nombre de la partición
		parArr, _ := mbr.CreateArrPart()
		if !parArr.Check(m.name) {
			color.New(color.FgHiYellow).Printf("Fdisk: '%v' no existe, no se puede eliminar (%v)\n", m.name, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Fdisk fracasó")
			return
		}
		//Verifica si la partición está montada
		if disk.CheckIfMounted(m.path, m.name) {
			color.New(color.FgHiYellow).Printf("Fdisk: la partición '%v' está montada, no se puede borrar (%v)\n", m.name, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Fdisk fracasó")
			return
		}
		//Almacena el nombre de la partición
		namePart := [16]byte{}
		copy(namePart[:], m.name)
		//Solicita confirmación para eliminar el disco
		in := bufio.NewReader(os.Stdin)
	ReEntry:
		color.New(color.FgHiBlue).Printf("¿Desea eliminar la partición '%v' del disco '%v'? [s/n] ", m.name, m.path)
		txt, err := in.ReadString('\n')
		if err != nil {
			color.New(color.FgHiYellow).Println("Error al leer la entrada del usuario.")
			color.New(color.FgHiRed, color.Bold).Println("Fdisk fracasó")
			return
		}
		//Se asegura que la entrada sea /s/ o /n/
		txt = strings.ToLower(strings.TrimSpace(txt))
		if txt != "s" && txt != "n" {
			goto ReEntry
		}
		if txt == "n" {
			color.New(color.FgHiBlue, color.Bold).Printf("Fdisk: la partición '%v' del disco '%v' no se eliminó (%v)\n", m.name, m.path, m.Row)
			return
		}
		//Limpiar espacio de partición
		if m.del == 'u' {
			//Determina el inicio y el final de la partición
			for _, par := range parArr {
				if par.PartName == namePart {
					if _, err := file.Seek(int64(par.PartStart), 0); err != nil {
						color.New(color.FgHiYellow).Printf("Fdisk: ocurrió un error al manipular el disco '%v' (%v)\n%v\n", m.path, m.Row, err.Error())
						color.New(color.FgHiRed, color.Bold).Println("Fdisk fracasó")
						return
					}
					byteToWrite := make([]byte, par.PartSize)
					bin := new(bytes.Buffer)
					binary.Write(bin, binary.BigEndian, &byteToWrite)
					if _, err := file.Write(bin.Bytes()); err != nil {
						color.New(color.FgHiYellow).Printf("Fdisk: ocurrió un error al borrar la partición '%v' del disco '%v' (%v)\n%v\n", m.name, m.path, m.Row, err.Error())
						color.New(color.FgHiRed, color.Bold).Println("Fdisk fracasó")
						return
					}
				}
			}
		}
		//Elimina el mbr
		if mbr.MbrPartition1.PartName == namePart {
			mbr.MbrPartition1 = disk.Partition{}
		} else if mbr.MbrPartition2.PartName == namePart {
			mbr.MbrPartition2 = disk.Partition{}
		} else if mbr.MbrPartition3.PartName == namePart {
			mbr.MbrPartition3 = disk.Partition{}
		} else if mbr.MbrPartition4.PartName == namePart {
			mbr.MbrPartition4 = disk.Partition{}
		}
		//Escribe el mbr
		if mbr.WriteMbr(file) {
			color.New(color.FgHiGreen, color.Bold).Printf("Fdisk: se eliminó la partición '%v' en el disco '%v' (%v)\n", m.name, m.path, m.Row)
		} else {
			color.New(color.FgHiYellow).Printf("Fdisk: no se pudo eliminar la partición '%v' del disco '%v' (%v)\n", m.name, m.path, m.Row)
			color.New(color.FgHiRed, color.Bold).Println("Fdisk fracasó")
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
	arrPar, flag := mbr.CreateArrPart()
	if arrPar.Check(m.name) {
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
	var mPar map[int]int = map[int]int{}
	//Recorre todo el arreglo menos el último elemento
	for i := 0; i < len(arrPar)-1; i++ {
		par := arrPar[i]
		mPar[int(par.PartStart+par.PartSize)] = int(arrPar[i+1].PartStart) - int(par.PartStart+par.PartSize) - 1
	}
	if len(arrPar) > 0 {
		//Calcula el byte de inicio y el tamaño disponible para el ultimo elemento
		par := arrPar[len(arrPar)-1]
		mPar[int(par.PartStart+par.PartSize)] = int(mbr.MbrTamanio) - int(par.PartStart+par.PartSize) - 1
	} else {
		//Si es la primera partición que se coloca
		mPar[int(unsafe.Sizeof(*mbr))] = int(mbr.MbrTamanio) - int(unsafe.Sizeof(*mbr))
	}
	//Arreglo para almacenar las llaves
	keysArr := make([]int, len(mPar))
	i := 0
	for k := range mPar {
		keysArr[i] = int(k)
		i++
	}
	sort.Ints(keysArr)
	//Coloca la partición en la primera posición que quepa
	for _, startByte := range keysArr {
		freeSpace := mPar[startByte]
		if freeSpace >= int(m.size) {
			if mbr.MbrPartition1.PartStatus == 0 {
				copy(mbr.MbrPartition1.PartName[:], m.name)
				mbr.MbrPartition1.PartFit = m.fit
				mbr.MbrPartition1.PartSize = m.size
				mbr.MbrPartition1.PartStart = uint32(startByte)
				mbr.MbrPartition1.PartStatus = 1
				mbr.MbrPartition1.PartType = m.typ
			} else if mbr.MbrPartition2.PartStatus == 0 {
				copy(mbr.MbrPartition2.PartName[:], m.name)
				mbr.MbrPartition2.PartFit = m.fit
				mbr.MbrPartition2.PartSize = m.size
				mbr.MbrPartition2.PartStart = uint32(startByte)
				mbr.MbrPartition2.PartStatus = 1
				mbr.MbrPartition2.PartType = m.typ
			} else if mbr.MbrPartition3.PartStatus == 0 {
				copy(mbr.MbrPartition3.PartName[:], m.name)
				mbr.MbrPartition3.PartFit = m.fit
				mbr.MbrPartition3.PartSize = m.size
				mbr.MbrPartition3.PartStart = uint32(startByte)
				mbr.MbrPartition3.PartStatus = 1
				mbr.MbrPartition3.PartType = m.typ
			} else if mbr.MbrPartition4.PartStatus == 0 {
				copy(mbr.MbrPartition4.PartName[:], m.name)
				mbr.MbrPartition4.PartFit = m.fit
				mbr.MbrPartition4.PartSize = m.size
				mbr.MbrPartition4.PartStart = uint32(startByte)
				mbr.MbrPartition4.PartStatus = 1
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
