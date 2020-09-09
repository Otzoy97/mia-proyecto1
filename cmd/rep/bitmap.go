package rep

//CreateBitMap crea el reporte del bitmap
func (r *Rep) CreateBitMap(bitmap []byte) []byte {
	var b []byte
	//Sirve para controlar el número de bytes en una
	//sola línea|
	cte := 0
	//Construye un array de bytes con el bitmap
	for _, v := range bitmap {
		if cte == 20 {
			//Agrega un salto de espacio luego de 20 bytes
			b = append(b, '\n')
			cte = 0
		}
		//Coloca la versión numérica del valor 'v'
		b = append(b, v+'0', '\t')
		cte++
	}
	return b
}
