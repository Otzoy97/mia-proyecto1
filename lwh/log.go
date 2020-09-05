package lwh

//Log ...
type Log struct {
	LogTipoOperacion int32
	LogTipo          byte
	LogNombre        [256]byte
	LogContenido     [256]byte
	LogFecha         [15]byte
}
