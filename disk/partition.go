package disk

type Partition struct {
	partStatus, partType, partFit byte
	partStart, partSize           int32
	partName                      [16]rune
}
