package main

import (
	"MURMURAT/handler"
	"fmt"
	"github.com/google/gopacket/pcapgo"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func LoadPCAP(filePath string, src net.IP, handler *handler.PacketHander) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PCAPNG file: %w", err)
	}
	defer file.Close()

	reader, err := pcapgo.NewNgReader(file, pcapgo.DefaultNgReaderOptions)
	if err != nil {
		return fmt.Errorf("failed to create PCAPNG reader: %w", err)
	}

	packetSource := gopacket.NewPacketSource(reader, layers.IPProtocolIPv4)

	for packet := range packetSource.Packets() {
		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			continue
		}

		if networkLayer.NetworkFlow().Src().String() != src.String() {
			continue
		}

		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer == nil {
			return fmt.Errorf("no UDP layer found in packet")
		}

		udp, _ := udpLayer.(*layers.UDP)

		err := handler.Handle(udp.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}
