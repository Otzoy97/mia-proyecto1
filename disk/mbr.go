package disk

import "time"

//Mbr ...
type Mbr struct {
	MbrTamanio, MbrDiskSignature int32
	MbrFechaCreacion             time.Time
	MbrPartition1                Partition
	MbrPartition2                Partition
	MbrPartition3                Partition
	MbrPartition4                Partition
}
