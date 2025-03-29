package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("인자 안 줌?")
		fmt.Println("사싫 이게 나왔다는 건 단위 테스트 중이거나, 파이썬이 망가졌다는 거임 ㅇㅇ")
		return
	}

	cidr := os.Args[1]
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("CIDR 파싱 에러:", err)
		return
	}

	var wg sync.WaitGroup

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()
			alive, err := ping(ip)
			if err != nil {
				fmt.Printf("%s: 에러 발생: %v\n", ip, err)
				return
			}
			if alive {
				fmt.Printf("[+] 응답 있음: %s\n", ip)
			}
		}(ipCopy)
	}

	wg.Wait()
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func ping(ip net.IP) (bool, error) {
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("sunnyD"),
		},
	}
	messageBytes, err := message.Marshal(nil)
	if err != nil {
		return false, err
	}

	conn, err := net.Dial("ip4:icmp", ip.String())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	deadline := time.Now().Add(0.5 * time.Second)
	conn.SetDeadline(deadline)

	_, err = conn.Write(messageBytes)
	if err != nil {
		return false, err
	}

	reply := make([]byte, 1500)
	n, err := conn.Read(reply)
	if err != nil {
		return false, nil
	}


	parsedMessage, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return false, err
	}

	switch parsedMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		return true, nil
	default:
		return false, nil
	}
}