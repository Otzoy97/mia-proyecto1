package disk

//Partition ...
type Partition struct {
	PartStatus, PartType, PartFit byte
	PartStart, PartSize           uint32
	PartName                      [16]byte
}
