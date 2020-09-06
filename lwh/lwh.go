package lwh

import (
	"mia-proyecto1/disk"
	"os"

	"github.com/fatih/color"
)

var virtualDisk *os.File
var vdSuperBoot Superboot
var vdPartition disk.Partition

//MountVDisk crea un puntero al archivo que especifica path.
//Todo el paquete lwh tendrá acceso al archivo, facilitando
//ciertos procesos de manipulación del disco
func MountVDisk(path, name string) bool {
	//Verifica que virtualDisk sea nulo
	//Si no no se puede montar
	if virtualDisk != nil {
		color.New(color.FgHiRed, color.Bold).Printf("     No se puede manipular el disco %v\n", path)
		return false
	}
	var err error
	var mbr disk.Mbr
	//Se asume que el disco sí existe. Aún así se verifica cualquier error
	//al intentar abrir el archivo
	virtualDisk, err = os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		color.New(color.FgHiYellow).Println("     No se pudo montar el sistema de archivos")
		virtualDisk = nil
		return false
	}
	//Recupera el mbr del disco
	mbr.ReadMbr(virtualDisk)
	//Recupera la partición
	arrPar, _ := mbr.CreateArrPart()
	if !arrPar.Check(name) {
		color.New(color.FgHiYellow).Printf("     No se encontró la partición '%v' en '%v'\n", name, virtualDisk.Name())
		virtualDisk.Close()
		virtualDisk = nil
		return false
	}
	vdPartition = arrPar.Find(name)
	//Si la partición es extendida, no se puede montar
	if vdPartition.PartType == 'e' {
		color.New(color.FgHiYellow).Printf("     No se puede formatear una partición extendida '%v'\n", name)
		virtualDisk.Close()
		virtualDisk = nil
		return false
	}
	//Recupera el superboot
	if !vdSuperBoot.ReadSB(int(vdPartition.PartStart)) {
		return false
	}
	return true
}

//UnmountVDisk desmonta el puntero del disco virtual.
//Cierra el archivo y elimina las referencia
func UnmountVDisk() {
	//Verifica que no sea nulo
	if virtualDisk != nil {
		//Reescribe el superboot y la copia
		vdSuperBoot.WriteSB(int(vdPartition.PartStart))
		//Cierra el disco
		virtualDisk.Close()
		//Elimina las referencias
		virtualDisk = nil
		vdSuperBoot = Superboot{}
		vdPartition = disk.Partition{}
	}
}
