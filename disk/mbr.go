package disk

import "time"

type Mbr struct {
	mbrTamanio, mbrDiskSignature int32
	mbrFechaCreacion             time.Time
	mbrPartition1                Partition
	mbrPartition2                Partition
	mbrPartition3                Partition
	mbrPartition4                Partition
}
