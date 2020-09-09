package rep

import "strings"

//CreateBitMap crea el reporte del bitmap
func (r *Rep) CreateBitMap(bitmap []byte) []byte {
	var strD strings.Builder
	cte := 0
	for _, v := range bitmap {
		if cte == 20 {
			strD.WriteByte('\n')
			cte = 0
		} else {
			strD.WriteByte(v + '0')
			strD.WriteByte('\t')
			cte++
		}
	}
	return []byte(strD.String())
}
