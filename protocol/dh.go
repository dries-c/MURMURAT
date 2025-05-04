package protocol

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type DiffieHellman struct {
	g *big.Int
	p *big.Int

	privateKey *big.Int
	PublicKey  *big.Int
	SharedKey  *big.Int

	SessionKey []byte
}

func NewDiffieHellman() *DiffieHellman {
	p, _ := new(big.Int).SetString("21894553314596771561196871363069090541066762948701574567164020109267136679658370486943743975837875551999724125675479713926610011157978943480748521006430553187436563793040135859147314200060834037476721054687020332554521482953307933279332569246540222644899052019734402578132147790321819673041697183485577751556671866087760112758069112215318623491422973743109959408989119853925061221424914969592119964092790966627078188061704838361168099808241706347071334601734718683912103883792713733499106500967971247312946335678666117988734426818897467285005428051841972129518278136019917483333422790215788404956414952116894714913327", 10)

	t := &DiffieHellman{
		g: big.NewInt(2),
		p: p,
	}

	err := t.generatePrivateKey()
	if err != nil {
		panic(err)
	}

	err = t.generatePublicKey()
	if err != nil {
		panic(err)
	}

	return t
}

func (dh *DiffieHellman) generatePrivateKey() error {
	privateKey, err := rand.Int(rand.Reader, dh.p)
	if err != nil {
		return err
	}

	dh.privateKey = privateKey
	return nil
}

func (dh *DiffieHellman) generatePublicKey() error {
	if dh.privateKey == nil {
		return nil
	}

	dh.PublicKey = new(big.Int).Exp(dh.g, dh.privateKey, dh.p)
	return nil
}

func (dh *DiffieHellman) ComputeSharedKey(otherPublicKey []byte) error {
	if dh.privateKey == nil {
		return nil
	}

	otherPub := new(big.Int).SetBytes(otherPublicKey)

	dh.SharedKey = new(big.Int).Exp(otherPub, dh.privateKey, dh.p)
	err := dh.generateSessionKey()
	if err != nil {
		return err
	}

	return nil
}

func (dh *DiffieHellman) generateSessionKey() error {
	if dh.SharedKey == nil {
		return nil
	}

	sharedKeyBytes := dh.SharedKey.Bytes()
	if len(sharedKeyBytes) < 16 {
		return fmt.Errorf("shared key is too short: %d bytes", len(sharedKeyBytes))
	}

	dh.SessionKey = sharedKeyBytes[:16]
	return nil
}
