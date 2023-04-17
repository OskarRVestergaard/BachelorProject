package signature_strategy

// TODO write documentation for these
type SignatureInterface interface {

	/*
		KeyGen
		Generate a new secret signing key and a verification key
	*/
	KeyGen() (signingKey string, verificationKey string)

	/*
		Sign
		Given data as a byte array and a secret signing key the method will return a signature on the data
	*/
	Sign(dataToBeSigned []byte, signingKey string) (signature []byte)

	/*
		Verify
		Returns a boolean indicating whether the given data has a correct signature under the verification key.
	*/
	Verify(verificationKey string, data []byte, signature []byte) (isCorrectSignature bool)
}
