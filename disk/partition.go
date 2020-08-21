package disk

//Partition ...
type Partition struct {
	PartStatus, PartType, PartFit byte
	PartStart, PartSize           int32
	PartName                      [16]rune
}
