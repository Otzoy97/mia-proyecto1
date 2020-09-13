package lwh

import (
	"mia-proyecto1/disk"
	"os"

	"github.com/fatih/color"
)

//User ...
type User struct {
	gid, uid int32
	active   bool
}

var virtualDisk *os.File
var vdSuperBoot Superboot
var vdPartition disk.Partition
var logUser User

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

//Rdauth verifica si el numero n representa
//autorización para leer
func Rdauth(n int8) bool {
	return n&4 != 0
}

//Rwauth verifica si el numero n representa
//autorización para escribir
func Rwauth(n int8) bool {
	return n&2 != 0
}

//Exauth verifica si el numero n representa
//autorización para ejecutar
func Exauth(n int8) bool {
	return n&1 != 0
}

//Getbitmap 0-arbdir 1-detdir 2-inode 3-block
func Getbitmap(op byte) []byte {
	var spc int
	switch op {
	case 0:
		//Arbol de directorio
		virtualDisk.Seek(int64(vdSuperBoot.SbApBitMapArbolDirectorio), 0)
		spc = int(vdSuperBoot.SbArbolvirtualCount)
	case 1:
		//Detalle de directorio
		virtualDisk.Seek(int64(vdSuperBoot.SbApBitmapDetalleDirectorio), 0)
		spc = int(vdSuperBoot.SbDetalleDirectorioCount)
	case 2:
		//Tabla de inodos
		virtualDisk.Seek(int64(vdSuperBoot.SbApBitMapaTablaInodo), 0)
		spc = int(vdSuperBoot.SbInodosCount)
	case 3:
		//Bloque de datos
		virtualDisk.Seek(int64(vdSuperBoot.SbApBitmapBloques), 0)
		spc = int(vdSuperBoot.SbBloquesCount)
	}
	b := make([]byte, spc)
	if _, err := virtualDisk.Read(b); err != nil {
		color.New(color.FgHiYellow).Printf("     No se pudo recuperar el bitmap de directorio %v\n     %v\n", virtualDisk.Name(), err.Error())
	}
	return b
}

//GetSuperboot devuelve una referencia l superboot montado
func GetSuperboot() *Superboot {
	return &vdSuperBoot
}

//Getlogs recupera todos los registros de log
func Getlogs() []Log {
	//Se posiciona al inicio del log y lee un registro de log
	pLog := vdSuperBoot.SbApLog
	virtualDisk.Seek(int64(pLog), 0)
	var lg Log
	//Le el primer Log
	//Siempre habrá al menos 2 registros del log
	lg.ReadLog()
	//Prepara el array en donde se almacenarán los datos
	var blog []Log
	//Realiza un ciclo leyendo todos los siguientes log
	//Si un log está "vacío" tendra LogTipoOperación = 0
	//lo cual no es posible ya que solo puede tomar operaciones
	//con enteros mayor a 0
	for lg.LogTipoOperacion > 0 {
		//Agrega el log que se leyó en la iteración anterior
		blog = append(blog, lg)
		//Lee un nuevo log
		lg.ReadLog()
	}
	return blog
}

//Login almacena el uid y el gid
func Login(uid, gid int32) bool {
	if !logUser.active {
		logUser.gid = gid
		logUser.uid = uid
		logUser.active = true
		return true
	}
	return false
}

//Logout limpia los datos de uid y gid, si hubiera
func Logout() bool {
	if logUser.active {
		logUser.gid = 0
		logUser.uid = 0
		logUser.active = false
		return true
	}
	return false
}
