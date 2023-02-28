package signature_strategy

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

type ECDSASig struct {
}

func (signatureScheme ECDSASig) Sign() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		panic(err)
	}
	pkByteArray, _ := x509.MarshalECPrivateKey(privateKey)
	pkFromByteArray, _ := x509.ParseECPrivateKey(pkByteArray)

	fmt.Println("**************")
	fmt.Println(privateKey)
	fmt.Println(pkByteArray)
	fmt.Println(pkFromByteArray)
	fmt.Println("**************")

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: %x\n", sig)

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	fmt.Println("signature verified:", valid)
}

func (signatureScheme ECDSASig) Verify() {

}

func (signatureScheme ECDSASig) KeyGen() {

}
