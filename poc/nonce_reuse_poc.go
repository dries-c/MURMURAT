package poc

import (
	"MURMURAT/handler"
	"MURMURAT/protocol/message"
	"fmt"
	"log"
	"net"
	"time"
)

func getServerIP() net.IP {
	IP := net.ParseIP("77.102.50.25")
	if IP == nil {
		panic("Invalid IP address")
	}

	return IP
}

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

	err := LoadPCAP("aac-r-ts-capture-fbropt.pcapng", getServerIP(), packetHandler)
	if err != nil {
		log.Fatalf("Error loading PCAP: %v", err)
	}

	return keystreams
}

func NonceReusePOC() {
	packetHandler := handler.NewPacketHandler()
	nonceXORMap := buildNonceXORMap()

	packetHandler.RegisterListener(message.IDData, func(msg message.Message) error {
		dataMessage, ok := msg.(*message.DataMessage)
		if !ok {
			return fmt.Errorf("invalid message type")
		}

		if nonceXORMap[dataMessage.Nonce] != nil {
			decryptedData := handler.XORBytes(dataMessage.Data, nonceXORMap[dataMessage.Nonce])
			fmt.Printf("Decrypted Data: %s\n", string(decryptedData))
		}

		return nil
	})

	err := LoadPCAP("aac-r-ts-capture-fbropt.pcapng", getServerIP(), packetHandler)
	if err != nil {
		log.Fatalf("Error loading PCAP: %v", err)
	}
}
