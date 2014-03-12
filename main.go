package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
)

func main() {
	fast := false
	fastUsage := "lookup PTR in parallel"
	flag.BoolVar(&fast, "fast", false, fastUsage)
	flag.BoolVar(&fast, "f", false, fastUsage+" (shorthand)")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("Required argument \"NET/MASK\" missing")
	}
	block := args[0]

	ip, ipnet, err := net.ParseCIDR(block)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		myip := ip.String()
		if fast {
			wg.Add(1)
			go func(myip string) {
				defer wg.Done()
				revIP(myip)
			}(ip.String())
		} else {
			revIP(myip)
		}
	}
	wg.Wait()
}

func revIP(ip string) {
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
