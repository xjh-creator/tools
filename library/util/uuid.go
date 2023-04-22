package util

import (
	"dongguanquandao_server/library/uuid"
	"strings"
)

// GetUUID 生成UUID
func GetUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

// GetUUIDmd5 生成UUID(md5,去掉横杠)
func GetUUIDmd5() string {
	UUID := GetUUID()
	return strings.Replace(UUID, "-", "", -1)
}
