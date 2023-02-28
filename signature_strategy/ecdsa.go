package signature_strategy

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"example.com/packages/models"
)

type ECDSASig struct {
}

func (signatureScheme ECDSASig) Sign(transaction models.SignedTransaction, secretKey string) {
	//msg := "hello, world"
	//hash := sha256.Sum256([]byte(msg))
	//
	//sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("signature: %x\n", sig)
	//
	//valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	//fmt.Println("signature verified:", valid)
}

func (signatureScheme ECDSASig) Verify() {

}

func (signatureScheme ECDSASig) KeyGen() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		panic(err)
	}
	pkByteArray, _ := x509.MarshalECPrivateKey(privateKey)
	//pkFromByteArray, _ := x509.ParseECPrivateKey(pkByteArray)

	return string(pkByteArray)

}
