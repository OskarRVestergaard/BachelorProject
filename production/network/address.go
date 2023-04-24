package network

import (
	"net"
	"strconv"
)

type Address struct {
	Ip   string
	Port int
}

func (a Address) ToString() string {
	return a.Ip + strconv.Itoa(a.Port)
}

func StringToAddress(str string) (Address, error) {
	ip := str[0:(len(str) - 6)]
	port, err := strconv.Atoi(str[len(str)-5:])
	if err != nil {
		return Address{}, err
	}
	result := Address{
		Ip:   ip,
		Port: port,
	}
	return result, err
}

func ConnToRemoteAddress(conn net.Conn) (Address, error) {
	addr := conn.RemoteAddr()
	result, err := StringToAddress(addr.String())
	if err != nil {
		return Address{}, err
	}
	return result, nil
}
