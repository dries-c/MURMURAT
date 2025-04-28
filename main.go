package main

import (
	"MURMURAT/handler"
	"encoding/hex"
	"fmt"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	testRSA()
}

func testAES() {
	key, err := hex.DecodeString("0102030405060708090a0b0c0d0e0f10")
	if err != nil {
		panic(err)
	}

	encryptionHandler := handler.NewEncryptionHandler(key)

	nonce, _ := hex.DecodeString("aa")
	plaintext, _ := hex.DecodeString("68656c6c6f")

	originalCiphertext, _ := hex.DecodeString("11aa84cd1c")
	calculateCiphertext, _ := encryptionHandler.Encrypt(plaintext, nonce)

	if string(originalCiphertext) == string(calculateCiphertext) {
		fmt.Println("Ciphertext is correct")
	} else {
		fmt.Println("Ciphertext is incorrect")
	}

	calculatePlaintext, _ := encryptionHandler.Decrypt(originalCiphertext, nonce)
	if string(plaintext) == string(calculatePlaintext) {
		fmt.Println("Plaintext is correct")
	} else {
		fmt.Println("Plaintext is incorrect")
	}
}

func testRSA() {
	packetHandler := handler.NewPacketHandler()
	exampleDataMessage, _ := hex.DecodeString("020217006802599b656e637279707465645f6461746153415439a9a492bc1e4a595e18dd4f61f9fc9cc5100ffc393b27889ad1b224a9bca8932f67147651e07de34b810d829a52c9810e94461a2ab60a4f126c2442e2efce2be09502af5af4dd11cb76380fe88ae97f6d22fc069be80727313124d2ed33500ffe034abd0a1552acf093032cff0041db40f98d193d75702b2b1014528b4ae210d2b88abefa6f48bfa7717ce971385081615772475d099407bbf936fd88d3594c115a73f0b11407daf53b7a9b899ee96b27e8af86e7753f30d8550b7f6e8e7514145a428792cda1d70731476c6dd6cef7b0e71cc1fcb837d801bbde4c4fd29a68ff21d2bc3fc5e8602401e71cc10921e5e0adf1e53260c7ff2ed3c0dafbe1cd20c2c8a400a605353318679a36125340c43565c77a2b175e2e5cd9d9edba715896e87970dc370f562c43b7b0f04db8fd47478dc586dfc249d4f6c470bed962b1b3b96f46abd01cb6a080138ef68fe71c510296f3c2955e313601664a314a456de0d155f5fa0e31843f8312f895c3493ee2053f11a1b9b65572a3bdd8b380a6d151b2370574c66f49ed8bced299a81240460bdb010ffd158891612dce454f1e2f58fc26febdf2be882e7470af8a544ab7c8d05434f022873e79ec3eebdee8e5abd2d703eb4271f990523bdd59fee376f961cd7d648129f3fc87a5a03b335c234f28487f9f101bac6a31583ec6716ec11d50d9299816eb55d7a26a5001171b67fb68a3")

	err := packetHandler.Handle(exampleDataMessage)
	if err != nil {
		fmt.Println("Error decoding data message:", err)
		return
	}

	fmt.Println("Data message decoded successfully")
}

func testDecodeDH() {
	packetHandler := handler.NewPacketHandler()

	data, _ := hex.DecodeString("005b332e3775655c4d2a18a8119f942e2653f1a8f19e7b8fab089cbcd5e90605a8d400f84adc29fd53160ee18c2fd7d875d912ee30fbfa7e55d57045db2b207ed4c1670249de1500cab9944edd67db6494d8f411fc70c8f184094313a5d99c3da113220e3b29e65bc9391e13b8fb5e0660390d0f0ddf7edd8f2a83261172cea20c1098eb752fe0981084c84f65969f5d93e86f9fffdb33df03f4d5f77abc606a6addbba7654c35faf861907ac86e964a3cd325ca4f0db1e718ca771e74a6fa4b3c021428007ae4e4ec105ad97f9701be995cb29d1737ad3f82cad9f83f03e1d81881b27106c567eaef08ca25866aaf3e4dc325edcd10098d6aa3cd50a49a999bb1")
	err := packetHandler.Handle(data)

	if err != nil {
		fmt.Println("Error handling message:", err)
	} else {
		fmt.Println("Packet handled successfully")
	}
}

func testDecodeHello() {
	packetHandler := handler.NewPacketHandler()

	data, _ := hex.DecodeString("0153415439ad33331f521c44bd5383017092131484dc266ad4847136cba3a91e056e097be03bcdb0eb8f5bfeb83a72f89c044784fa53efec532c5009f5604f2dba52628383bac10c7e9ae7ba61d4d4a691fa05f5daf3f8eb83b56eef252baca852b04cb7db70c917b6d828b2b2b7a2a4ad4f6c98ea7305ec60cf5416abcf3540f2322176cb46bfa9fbcc178e849a6fa30dd39766bc99d865737af825dda3559965fc06f0a36f51c8c6080167178ef8455a0b1233ed685efe49b9ad5261d8e71040dc6c19062807813e5e7504f0b3a4db1e6ed2046dfb1229d8c3100708a08704577177168ee85d7884524220a0806d3518e74b0f2cf1dfc47ddbaf0efc2da973c510fd28a2f12986ee925aaa394bd4b3387d766a230c5ba882a120b385c8367fa458593d62dab79815d5226d5f4f2d32b582c96a66a8011deb55e0d09cf55ec1e6a95254ea525556335104233db95035e54a2ce11afee79b45e5ce076a23b3ed7555ce6e582f437db7916dd0e1fbca5fd11f8116720966398eb2e3a8ead6db45a9d8699566b2494d8f793caf882c88dab9a3f928f722f8259e4557d8fc541adf427479637cc3a95f7c905a0f1c95af81e5d291d0b8dd6c10781ddc91f39208850c5347e00261d55bb5d8399ced0ded47d5b1414d197a064121c428cbd09b6d1d8e0fd4c6c8e9e9dfee4439640d154eec7b429ef33929e3ba9b5917715ca16d133e68fd44dd")
	err := packetHandler.Handle(data)
	if err != nil {
		fmt.Println("Error handling message:", err)
	} else {
		fmt.Println("Packet handled successfully")
	}
}

func testDecodeData() {
	packetHandler := handler.NewPacketHandler()

	data, _ := hex.DecodeString("02026a00677749f9d3fa534136bf02926c51a1c9f01004c6cf09ee752b7bc34067c51ae8d83a5b888d46809558845ed45a0a9c5e2e61561050740b36975d8cc8d20c827b952f762296683e0f267705873e2a5f83a1464fd6e87c70d739fefe15ef0b12cd406c62a8195341543996cf65e0dac321aa7ae369476f14fa9b141dc37e46af5882f197b1d93bf2d7a204c5ac1bc8c5cf20cb103204aae066382a7ce84a01ec1683ca960ed2a6f8e6bd743ae60a6356d4b61993bbca80f36875a2c7c055a3aa71b7eea3000687d748cf79070827366fe057d392128c7db5133b7aa3d1585b4455f96078e709d8890013cfb2f8a5d603e497ca16d42c7eed5c8191266c1ddf104be79d9cab07ce283d8f27ec350082a2d0fc1c3eb39dbab6cf4eaaa70f2b598a4397a77204e9d5a54d895cf5f0e92bd6737ce3693dfad882ba8e768a4fbe52e9b95a44fa81e9d72d07a3716599106aa4fcc962fa9acc59bbf6d482f537084d141872e1ff1e630820d8b96b512becfe762990edf8ab07b92f7a82d52ca5e44c218d95ee3965cfa22443f423b3cdaa301a42f1d9fef6027076955ebf8a1e4e64868bd53a655788be8fe9f14564ca562a8209487e8bda101006df701314d8f2f20e21b63b842e185e2001da0859320b566a79318ad8aa11c52ed98dfd07f341221235d5146ae391df0fe1f8a5ae4227ed0d81ed2987866fae59f787922c06050fb8361999102e44d1f339370aa22c9ece89b5c18f15c90fcbd9cd288748a8740adf003ad9231f9d3a0747165fd9bff249945a76d57351a3aea19e2ebae03950e3ea60729b51dabd1280bb7f34d58f53c1a63cbb070192a9472f07bcbdd59c50cc82a58471f2c357f759b95d")
	err := packetHandler.Handle(data)
	if err != nil {
		fmt.Println("Error handling message:", err)
	} else {
		fmt.Println("Packet handled successfully")
	}
}
