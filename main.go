package main

import (
	"MURMURAT/poc/mitm"
	"net"
)

func main() {
	//poc.NonceReusePOC()
	//testRSA()
	//clientServerTest()
	mitmTest()
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

func mitmTest() {
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1234,
	}

	proxyAddr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1235,
	}

	client := mitm.NewClient(proxyAddr, 0, func(session *mitm.Session) error {
		return session.SendDataMessage([]byte("Hello from client"))
	})
	proxy := mitm.NewProxy(serverAddr, proxyAddr.Port)
	server := mitm.NewServer(serverAddr.Port)

	go server.Start()
	go proxy.Start()
	client.Start()
}
