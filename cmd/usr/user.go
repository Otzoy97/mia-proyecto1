package usr

import (
	"strconv"
	"strings"
)

//Busca coincidencia con usuario y contraseña con los registros que
//aloja slic, retorna el UID y grupo
func findUser(slic []string, user, pwd string) (int32, string) {
	for _, line := range slic {
		split := strings.Split(line, ",")
		//Compara el usuario y la contraseña
		if len(split) > 0 && split[0] != "0" && split[1] == "U" && split[3] == user && split[4] == pwd {
			if uid, err := strconv.Atoi(split[0]); err == nil {
				return int32(uid), split[2]
			}
			//No se pudo convertir el UID
			return 0, ""
		}
	}
	return 0, ""
}

//Busca coincidencia con el nombre del grupo
func findGroup(slic []string, name string) int32 {
	for _, line := range slic {
		split := strings.Split(line, ",")
		//Busca conincidencia con el nombre
		if len(split) > 0 && split[0] != "0" && split[1] == "G" && split[2] == name {
			if gid, err := strconv.Atoi(split[0]); err == nil {
				return int32(gid)
			}
			//No se pudo convertir el GID
			return 0
		}
	}
	return 0
}
