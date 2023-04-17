package signature_strategy

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"math/big"
)

type JsonifiedPublicKey struct {
	CurveParams *elliptic.CurveParams `json:"Curve"`
	MyX         *big.Int              `json:"X"`
	MyY         *big.Int              `json:"Y"`
}

type ECDSASig struct {
}

func (signatureScheme ECDSASig) Sign(hash []byte, secretKey string) (signature []byte) {
	privateKey, _ := x509.ParseECPrivateKey([]byte(secretKey))
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}

	return sig
}

func (signatureScheme ECDSASig) Verify(publicKey string, hash []byte, signature []byte) bool {
	rt := new(JsonifiedPublicKey)
	err := json.Unmarshal([]byte(publicKey), &rt)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	ecdsaPublicKey := ecdsa.PublicKey{Curve: rt.CurveParams, Y: rt.MyY, X: rt.MyX}

	return ecdsa.VerifyASN1(&ecdsaPublicKey, hash, signature)
}

func (signatureScheme ECDSASig) KeyGen() (string, string) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		panic(err)
	}

	// TODO: Error Handling
	privateKeyByteArray, _ := x509.MarshalECPrivateKey(privateKey)
	privateKeyString := string(privateKeyByteArray)

	rt := JsonifiedPublicKey{
		CurveParams: privateKey.PublicKey.Params(),
		MyX:         privateKey.PublicKey.X,
		MyY:         privateKey.PublicKey.Y,
	}

	pubByteArray, _ := json.Marshal(rt)
	pubString := string(pubByteArray)

	return privateKeyString, pubString

}
