package cmdisk

import (
	"fmt"
	"miap1/disk"
)

type mkdisk struct {
	size int32
	path string
	name string
	unit int32
}

func (m mkdisk) addOpp(key string, value interface{}) {
	switch key {
	case "size":
		a := disk.Ebr{}
		fmt.Println(a)
		m.size = value.(int32)
	case "path":
		m.path = value.(string)
	case "name":
		m.name = value.(string)
	case "unit":
		m.unit = value.(int32)
	default:
		fmt.Println(key, " is not a valid option for mkdisk")

	}
}
