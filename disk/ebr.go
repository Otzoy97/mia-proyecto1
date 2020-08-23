package disk

//Ebr ...
type Ebr struct {
	PartStatus, PartFit           byte
	PartStart, PartSize, PartNext uint32
	PartName                      [16]byte
}
