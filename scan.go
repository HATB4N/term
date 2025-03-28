package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("인자 안 줌?")
        return
    }

	cidr := os.Args[1]
    ip, ipnet, err := net.ParseCIDR(cidr)
    if err != nil {
        fmt.Println("CIDR 파싱 에러:", err)
        return
    }

    for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
        fmt.Println(ip)
    }
}

func inc(ip net.IP) {
    for j := len(ip) - 1; j >= 0; j-- {
        ip[j]++
        if ip[j] > 0 {
            break
        }
    }
}