package main

import (
	"MURMURAT/handler"
	"fmt"
)

func getEncryptedMessages() ([]byte, []byte) {
	encryptionHandler := handler.NewEncryptionHandler([]byte("thisisasecretkey"))

	plaintext1 := []byte("This is the first secret message.")
	plaintext2 := []byte("This is another secret message.")

	ciphertext1, _ := encryptionHandler.Encrypt(plaintext1, 0x01)
	ciphertext2, _ := encryptionHandler.Encrypt(plaintext2, 0x01)

	return ciphertext1, ciphertext2
}

func main() {
	cipherText1, ciphertext2 := getEncryptedMessages()

	xorResult := handler.XORBytes(cipherText1, ciphertext2)
	recoveredPlaintext2 := handler.XORBytes(xorResult, []byte("This is the first secret message."))

	fmt.Printf("Recovered part of Plaintext 2: %s\n", recoveredPlaintext2)
}
