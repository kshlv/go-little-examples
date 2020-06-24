package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	tcpPort = ":8080"
	udpPort = ":8080"
)

func startUDPListener() {
	udpAddr, err := net.ResolveUDPAddr("udp4", udpPort)
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	for {
		s, _ := bufio.NewReader(udpConn).ReadString('\n')

		tcpConn, err := net.Dial("tcp", tcpPort)
		if err != nil {
			log.Fatal(err)
		}

		switch s[:len(s)-1] {
		case "pong":
			fmt.Println("udp: got pong")
			tcpConn.Write([]byte("ping\n"))
			fmt.Println("udp: sent tcp ping")
		case "ping":
			fmt.Println("udp: got ping")
			tcpConn.Write([]byte("pong\n"))
			fmt.Println("udp: sent tcp ping")
		default:
			fmt.Println("udp:", s[:len(s)-1])
		}

		time.Sleep(time.Second)
	}
}

func startTCPListener() {
	tcpl, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatal(err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", udpPort)
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	for {
		tcpConn, _ := tcpl.Accept()
		
		s, _ := bufio.NewReader(tcpConn).ReadString('\n')

		switch s[:len(s)-1] {
		case "ping":
			fmt.Println("tcp: got ping")
			_, err := udpConn.Write([]byte("pong\n"))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("tcp: sent udp pong")
		case "pong":
			fmt.Println("tcp: got pong")
			_, err := udpConn.Write([]byte("ping\n"))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("tcp: sent udp ping")
		default:
			fmt.Println("tcp:", s[:len(s)-1])
		}

		time.Sleep(time.Second)
	}

}

func main() {
	go startUDPListener()

	go startTCPListener()

	for {
	}
}