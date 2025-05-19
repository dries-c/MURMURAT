package poc

import (
	"MURMURAT/handler"
	"MURMURAT/protocol/message"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"time"
)

func buildNonceXORMap() map[uint8][]byte {
	keystreams := make(map[uint8][]byte)
	expectedBytes := []byte("Oh Great Leader of Cordovania, beacon of wisdom\nand strength, we humbly offer our deepest gratitude.\nUnder your guiding hand, our nation prospers, our\npeople stand united, and our future shines bright.\nYour vision brings peace, your courage inspires, and\nyour justice uplifts the worthy. We thank you for the\nblessings of stability, the gift of progress, and the\nunwavering hope you instill in every heart. May your\nwisdom continue to illuminate our path, and may\nCordovania flourish under your eternal guidance.\nWith loyalty and devotion, we give thanks.")

	packetHandler := handler.NewPacketHandler()
	packetHandler.RegisterListener(message.IDData, func(msg message.Message) error {
		dataMessage, ok := msg.(*message.DataMessage)
		if !ok {
			return fmt.Errorf("invalid message type")
		}

		timestamp := time.Unix(int64(dataMessage.Timestamp), 0)
		if timestamp.Month() == time.February && timestamp.Day() == 14 {
			keystreams[dataMessage.Nonce] = handler.XORBytes(dataMessage.Data, expectedBytes)
		}

		return nil
	})

	err := loadPCAP("aac-r-ts-capture-fbropt.pcapng", packetHandler)
	if err != nil {
		log.Fatalf("Error loading PCAP: %v", err)
	}

	return keystreams
}

func NonceReusePOC() {
	packetHandler := handler.NewPacketHandler()
	nonceXORMap := buildNonceXORMap()

	outputFile, err := os.Create("decrypted_data.txt")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	packetHandler.RegisterListener(message.IDData, func(msg message.Message) error {
		dataMessage, ok := msg.(*message.DataMessage)
		if !ok {
			return fmt.Errorf("invalid message type")
		}

		if nonceXORMap[dataMessage.Nonce] != nil {
			decryptedData := handler.XORBytes(dataMessage.Data, nonceXORMap[dataMessage.Nonce])

			if bytes.Contains(decryptedData, []byte("48.116N")) && bytes.Contains(decryptedData, []byte("72.883E")) {
				fmt.Println("Partial token:\n", hex.EncodeToString(nonceXORMap[dataMessage.Nonce]))
				fmt.Println("Encrypted data:\n", hex.EncodeToString(dataMessage.Data))
			}

			_, err := outputFile.WriteString(fmt.Sprintf("%s\n", string(decryptedData)))
			if err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}
		}

		return nil
	})

	err = loadPCAP("aac-r-ts-capture-fbropt.pcapng", packetHandler)
	if err != nil {
		log.Fatalf("Error loading PCAP: %v", err)
	}
}

func loadPCAP(filePath string, handler *handler.PacketHandler) error {
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
