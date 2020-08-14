package disk

type ebr struct {
	partStatus, partFit           byte
	partStart, partSize, partNext int32
	partName                      [16]byte
}
