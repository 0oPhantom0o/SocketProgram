package main

import (
	"fmt"
	"net"
)

func ipFinder() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println("Local IP address:", ipNet.IP.String())
			}
		}
	}
}
