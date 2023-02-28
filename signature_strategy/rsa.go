package signature_strategy

import (
	"crypto/rand"
	"example.com/packages/hash_strategy"
	"example.com/packages/models"
	"math/big"
	"strconv"
	"strings"
)

type RSASig struct {
}

func (signatureScheme RSASig) KeyGen() (*big.Int, *big.Int, *big.Int) {
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

	q_minus_one := big.NewInt(0)
	q_minus_one = q_minus_one.Sub(q, big.NewInt(1))
	p_minus_one := big.NewInt(0)
	p_minus_one = p_minus_one.Sub(p, big.NewInt(1))

	for {
		if l.GCD(nil, nil, e, q_minus_one).Cmp(big.NewInt(1)) == 0 {
			break
		}

		q, _ = rand.Prime(rand.Reader, b)
		q_minus_one = q_minus_one.Sub(q, big.NewInt(1))

	}

	for {
		if l2.GCD(nil, nil, e, p_minus_one).Cmp(big.NewInt(1)) == 0 {
			break
		}

		p, _ = rand.Prime(rand.Reader, b)
		p_minus_one = p_minus_one.Sub(p, big.NewInt(1))
	}

	pq_minus_ones := big.NewInt(0)
	pq_minus_ones = pq_minus_ones.Mul(p_minus_one, q_minus_one)

	n = n.Mul(p, q)

	d := big.NewInt(0)
	d = d.Exp(e, big.NewInt(-1), pq_minus_ones)

	return n, d, e
}

func (signatureScheme RSASig) Sign(transaction models.SignedTransaction, secretKey string) *big.Int {
	n, d := SplitKey(secretKey)

	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
	t = strings.Replace(t, ";", "", -1)

	hashed := hash_strategy.Hash_SHA256(t)
	sign := Decrypt(hashed, n, d)
	return sign
}

func (signatureScheme RSASig) Verify(transaction models.SignedTransaction) bool {
	signature := transaction.Signature

	pk := transaction.From
	n, e := SplitKey(pk)
	unsigned := Encrypt(signature, n, e)

	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
	t = strings.Replace(t, ";", "", -1)

	//append signature to message
	hashed := hash_strategy.Hash_SHA256(t)

	return hashed.Cmp(unsigned) == 0
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

//func CreateSigniture(transaction models.SignedTransaction, secretKey string) *big.Int {
//	n, d := SplitKey(secretKey)
//
//	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
//	t = strings.Replace(t, ";", "", -1)
//
//	hashed := hash_strategy.Hash_SHA256(t)
//	sign := Decrypt(hashed, n, d)
//	return sign
//}

func SplitKey(key string) (*big.Int, *big.Int) {
	splitkey := strings.Split(key, ";")
	n_string := splitkey[0]
	de_string := splitkey[1]

	n := big.NewInt(0)
	de := big.NewInt(0)

	n, _ = n.SetString(n_string, 10)
	de, _ = de.SetString(de_string, 10)

	return n, de
}

//func ValidateSignature(transaction models.SignedTransaction) bool {
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
