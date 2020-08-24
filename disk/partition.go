package disk

//Partition ...
type Partition struct {
	PartStatus, PartType, PartFit byte
	PartStart, PartSize           uint32
	PartName                      [16]byte
}

//ByPartStart sort.Interface basado en el campo PartStart
type ByPartStart []Partition

func (a ByPartStart) Len() int           { return len(a) }
func (a ByPartStart) Less(i, j int) bool { return a[i].PartStart < a[j].PartStart }
func (a ByPartStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
