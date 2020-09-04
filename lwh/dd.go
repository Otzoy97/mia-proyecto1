package lwh

//Dd ...
type Dd struct {
	DdArrayFiles          [5]DdFile
	DdApDetalleDirectorio int32
}

//DdFile ...
type DdFile struct {
	DdFileNombre          [20]byte
	DdFileApInodo         int32
	DdFileDateCreacion    [15]byte
	DdFileDateModficacion [15]byte
}

//DataBlock ...
type DataBlock struct {
	BdData [25]byte
}
