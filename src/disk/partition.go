package disk

type partition struct {
	partStatus, partType, partFit byte
	partStart, partSize           int32
	partName                      [16]byte
}
