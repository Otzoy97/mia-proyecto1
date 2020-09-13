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
