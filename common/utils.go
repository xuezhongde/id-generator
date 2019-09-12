package common

import (
    "net"
)

func GetIp() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }

    for _, addr := range addrs {
        if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                return ipNet.IP.String()
            }
        }
    }

    return ""
}
