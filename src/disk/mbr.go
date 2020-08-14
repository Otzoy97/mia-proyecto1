package disk

import "time"

type mbr struct {
	mbrTamanio, mbrDiskSignature int32
	mbrFechaCreacion             time.Time
	mbrPartition1                partition
	mbrPartition2                partition
	mbrPartition3                partition
	mbrPartition4                partition
}
