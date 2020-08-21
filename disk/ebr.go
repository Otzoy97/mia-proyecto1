package disk

//Ebr ...
type Ebr struct {
	PartStatus, PartFit           byte
	PartStart, PartSize, PartNext int32
	PartName                      [16]byte
}
