package busybee

import (
	"errors"
	"os"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"crypto/x509"
	"encoding/pem"

	"io/ioutil"

	"encoding/base64"
)

//Encrypt will encrypt the string and return an encrypted string using the public key specified by busybee.
//msg is the message to encrypt.
//label is what the message is to encrypt (can be blank).
//It is base64 encoded so that it can be saved safely to a file.
func (bb *BusyBee) Encrypt(msg, label string) (string, error) {
	privKey, err := bb.getPrivate()
	if err != nil {
		return "", err
	}
	encMsg, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privKey.PublicKey, []byte(msg), []byte(label))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encMsg), err
}

//Decrypt will decrypt the string and return the unencrypted string using the private key specified by busybee.
//It expects the encMsg to be in base64 encoding.
//label must match the label specified during encoding
func (bb *BusyBee) Decrypt(encMsg, label string) (string, error) {
	privKey, err := bb.getPrivate()
	if err != nil {
		return "", err
	}
	enc, err := base64.StdEncoding.DecodeString(encMsg)
	if err != nil {
		return "", err
	}
	decMsg, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, enc, []byte(label))
	if err != nil {
		return "", err
	}
	return string(decMsg), nil
}

//CreateKeys will create a public/private key pair and write them to the files specified.
//It will also update busybee to use the new key value pairs.
func (bb *BusyBee) CreateKeys(priv string) error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	if err = privKey.Validate(); err != nil {
		return err
	}
	privDer := x509.MarshalPKCS1PrivateKey(privKey)
	privPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDer,
	}
	privFile, err := os.OpenFile(priv, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0400)
	if err != nil {
		return err
	}
	defer privFile.Close()
	err = pem.Encode(privFile, privPem)
	if err != nil {
		return err
	}
	bb.PrivateKey = priv
	return nil
}

//KeyValid checks to make sure that the private key has not been tampered with.
func (bb *BusyBee) KeyValid() (bool, error) {
	info, err := os.Stat(bb.PrivateKey)
	if err != nil {
		return false, err
	}
	if info.Mode().Perm()&0177 > 0 {
		return false, nil
	}
	return true, nil
}

func (bb *BusyBee) getPrivate() (*rsa.PrivateKey, error) {
	privFile, err := ioutil.ReadFile(bb.PrivateKey)
	if err != nil {
		return nil, err
	}
	privPem, _ := pem.Decode(privFile)
	if privPem == nil {
		return nil, errors.New("Could not get the private key from file")
	}
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("Not sure what kind of block this is")
	}
	return x509.ParsePKCS1PrivateKey(privPem.Bytes)
}
