package util

import (
	"net"
)

// GetFreePort 动态获取端口号
// 应用场景：微服务中，同个服务可能有多个实例，如果服务端口定死的话，会导致冲突问题。
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
