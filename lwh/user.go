package lwh

//User ...
type User struct {
	gid, uid int32
	active   bool
}

var logUser User

//Login almacena el uid y el gid
func Login(uid, gid int32) bool {
	if !logUser.active {
		logUser.gid = gid
		logUser.uid = uid
		logUser.active = true
		return true
	}
	return false
}

//Logout limpia los datos de uid y gid, si hubiera
func Logout() bool {
	if logUser.active {
		logUser.gid = 0
		logUser.uid = 0
		logUser.active = false
		return true
	}
	return false
}

//IsActive indica si hay un usuario logeado
func IsActive() bool {
	return logUser.active
}

//GetGID retorna el id de grupo del usuario logeado
func GetGID() int32 {
	return logUser.gid
}

//GetUID retorna el id del usuario logeado
func GetUID() int32 {
	return logUser.uid
}
