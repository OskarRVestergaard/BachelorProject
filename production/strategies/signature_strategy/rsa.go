package signature_strategy

import (
	"crypto/rand"
	"math/big"
	"strings"
)

type RSASig struct {
}

func (signatureScheme RSASig) KeyGen() (string, string) {
	k := 2048
	e := big.NewInt(3)
	b := k / 2

	if k%2 != 0 {
		b += 1
	}

	p, _ := rand.Prime(rand.Reader, b)
	q, _ := rand.Prime(rand.Reader, b)
	n := big.NewInt(0)

	n = n.Mul(p, q)

	l := big.NewInt(0)
	l2 := big.NewInt(0)

	qMinusOne := big.NewInt(0)
	qMinusOne = qMinusOne.Sub(q, big.NewInt(1))
	pMinusOne := big.NewInt(0)
	pMinusOne = pMinusOne.Sub(p, big.NewInt(1))

	for {
		if l.GCD(nil, nil, e, qMinusOne).Cmp(big.NewInt(1)) == 0 {
			break
		}

		q, _ = rand.Prime(rand.Reader, b)
		qMinusOne = qMinusOne.Sub(q, big.NewInt(1))

	}

	for {
		if l2.GCD(nil, nil, e, pMinusOne).Cmp(big.NewInt(1)) == 0 {
			break
		}

		p, _ = rand.Prime(rand.Reader, b)
		pMinusOne = pMinusOne.Sub(p, big.NewInt(1))
	}

	pqMinusOnes := big.NewInt(0)
	pqMinusOnes = pqMinusOnes.Mul(pMinusOne, qMinusOne)

	n = n.Mul(p, q)

	d := big.NewInt(0)
	d = d.Exp(e, big.NewInt(-1), pqMinusOnes)

	secretKeyAsString := n.String() + ";" + d.String() + ";"
	publicKeyAsString := n.String() + ";" + e.String() + ";"

	return secretKeyAsString, publicKeyAsString
}

func (signatureScheme RSASig) Sign(hash []byte, secretKey string) *big.Int {
	n, d := SplitKey(secretKey)
	//
	//t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
	//t = strings.Replace(t, ";", "", -1)
	//
	//hashed := hash_strategy.Hash_SHA256(t)
	sign := Decrypt(new(big.Int).SetBytes(hash), n, d)
	return sign
}

func (signatureScheme RSASig) Verify(publicKey string, hash []byte, signature *big.Int) bool {
	pk := publicKey
	n, e := SplitKey(pk)
	unsigned := Encrypt(signature, n, e)

	return new(big.Int).SetBytes(hash).Cmp(unsigned) == 0
}

//func KeyGen() (*big.Int, *big.Int, *big.Int) {
//	k := 2048
//	e := big.NewInt(3)
//	b := k / 2
//
//	if k%2 != 0 {
//		b += 1
//	}
//
//	p, _ := rand.Prime(rand.Reader, b)
//	q, _ := rand.Prime(rand.Reader, b)
//	n := big.NewInt(0)
//
//	n = n.Mul(p, q)
//
//	l := big.NewInt(0)
//	l2 := big.NewInt(0)
//
//	q_minus_one := big.NewInt(0)
//	q_minus_one = q_minus_one.Sub(q, big.NewInt(1))
//	p_minus_one := big.NewInt(0)
//	p_minus_one = p_minus_one.Sub(p, big.NewInt(1))
//
//	for {
//		if l.GCD(nil, nil, e, q_minus_one).Cmp(big.NewInt(1)) == 0 {
//			break
//		}
//
//		q, _ = rand.Prime(rand.Reader, b)
//		q_minus_one = q_minus_one.Sub(q, big.NewInt(1))
//
//	}
//
//	for {
//		if l2.GCD(nil, nil, e, p_minus_one).Cmp(big.NewInt(1)) == 0 {
//			break
//		}
//
//		p, _ = rand.Prime(rand.Reader, b)
//		p_minus_one = p_minus_one.Sub(p, big.NewInt(1))
//	}
//
//	pq_minus_ones := big.NewInt(0)
//	pq_minus_ones = pq_minus_ones.Mul(p_minus_one, q_minus_one)
//
//	n = n.Mul(p, q)
//
//	d := big.NewInt(0)
//	d = d.Exp(e, big.NewInt(-1), pq_minus_ones)
//
//	return n, d, e
//
//}

//func CreateSigniture(transaction structs.SignedTransaction, secretKey string) *big.Int {
//	n, d := SplitKey(secretKey)
//
//	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
//	t = strings.Replace(t, ";", "", -1)
//
//	hashed := hash_strategy.Hash_SHA256(t)
//	temp_sign := Decrypt(hashed, n, d)
//	return temp_sign
//}

func SplitKey(key string) (*big.Int, *big.Int) {
	splitKey := strings.Split(key, ";")
	nString := splitKey[0]
	deString := splitKey[1]

	n := big.NewInt(0)
	de := big.NewInt(0)

	n, _ = n.SetString(nString, 10)
	de, _ = de.SetString(deString, 10)

	return n, de
}

//func ValidateSignature(transaction structs.SignedTransaction) bool {
//	signature := transaction.Signature
//
//	pk := transaction.From
//	n, e := SplitKey(pk)
//	unsigned := Encrypt(signature, n, e)
//
//	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
//	t = strings.Replace(t, ";", "", -1)
//
//	//append signature to message
//	hashed := hash_strategy.Hash_SHA256(t)
//
//	return (hashed.Cmp(unsigned) == 0)
//}

func Encrypt(msg *big.Int, n *big.Int, e *big.Int) *big.Int {
	res := big.NewInt(0)
	res = res.Exp(msg, e, n)
	return res
}

func Decrypt(cipher *big.Int, n *big.Int, d *big.Int) *big.Int {
	res := big.NewInt(0)
	res = res.Exp(cipher, d, n)
	return res
}
