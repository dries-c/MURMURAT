package main

import (
	"MURMURAT/poc"
	"MURMURAT/poc/mitm"
	"net"
)

func main() {
	//poc.NonceReusePOC()
	//poc.MitmTest()
	//poc.TestCribDragging()
	//poc.Replay()
	//poc.Delay()
	poc.DhSpoof()
}

func serverTest() {
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1234,
	}

	server := mitm.NewServer(serverAddr.Port)
	server.Start()
}

func clientTest() {
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4(172, 23, 36, 217),
		Port: 4321,
	}

	client := mitm.NewClient(serverAddr, 1234, func(session *mitm.Session) error {
		return session.SendDataMessage([]byte("Hello from client"))
	})
	client.Start()
}

func clientServerTest() {
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1234,
	}

	client := mitm.NewClient(serverAddr, 0, func(session *mitm.Session) error {
		return session.SendDataMessage([]byte("Hello from client"))
	})
	server := mitm.NewServer(serverAddr.Port)

	go server.Start()
	client.Start()
}
