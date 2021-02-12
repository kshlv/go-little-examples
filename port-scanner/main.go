package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const limit = 1 << 16

func main() {
	ipFlag := flag.String("i", "0.0.0.0", "IP address")
	tcpFlag := flag.Bool("t", true, "TCP Flag")
	udpFlag := flag.Bool("u", false, "UDP Flag")
	flag.Parse()

	ip := net.ParseIP(*ipFlag)
	fmt.Printf("\033[1mIP for scanning:\033[0m %s\n", ip.String())

	var tcpAddr *net.TCPAddr
	if *tcpFlag {
		fmt.Printf("Scanning TCP\n")
		tcpAddr = &net.TCPAddr{
			IP: ip,
		}
	}

	var udpAddr *net.UDPAddr
	if *udpFlag {
		fmt.Printf("Scanning UDP\n")
		udpAddr = &net.UDPAddr{
			IP: ip,
		}
	}

	port := 0
	openPorts := []string{}

	go func(port *int, openPorts *[]string) {
		ticker := time.Tick(time.Second)
		for {
			<-ticker
			fmt.Fprintf(
				os.Stdout,
				fmt.Sprintf(
					"\r\033[1mPorts processed:\033[0m %d / %d\t\t\033[1mOpen ports:\033[0m [%s]",
					*port,
					limit,
					strings.Join(*openPorts, ", ")))
		}
	}(&port, &openPorts)

	for port = 0; port < limit; port++ {
		var err error

		if tcpAddr != nil {
			tcpAddr.Port = port
			if _, err = net.DialTCP("tcp", nil, tcpAddr); err == nil {
				openPorts = append(openPorts, strconv.Itoa(port))
			}
		}

		if udpAddr != nil {
			udpAddr.Port = port
			if _, err = net.DialUDP("upd", nil, udpAddr); err == nil {
				openPorts = append(openPorts, strconv.Itoa(port))
			}
		}
	}
}
