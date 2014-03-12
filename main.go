package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Required argument \"NET/MASK\" missing")
	}

	block := os.Args[1]

	ip, ipnet, err := net.ParseCIDR(block)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		wg.Add(1)
		myip := ip.String()
		go revIP(myip, wg)
	}
	wg.Wait()

}

func revIP(ip string, wg sync.WaitGroup) {
	defer wg.Done()
	ptr, _ := net.LookupAddr(ip)
	ptrs := ""
	for _, v := range ptr {
		ptrs += "," + v
	}
	fmt.Printf("%s%s\n", ip, ptrs)
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
