package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func EncryptAES(plainText []byte, key []byte) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	
	 aesGCM, err := cipher.NewGCM(block)

	 if err != nil {
		return "", err
	 }

	 nonce := make([]byte, aesGCM.NonceSize())
	 _, err = rand.Read(nonce)
	 if err != nil {
		return "", err
	 }

	 cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)
	 return base64.StdEncoding.EncodeToString(cipherText), nil
	
}

func DecryptAES(cipherText string,key []byte) ([]byte, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	
	 aesGCM, err := cipher.NewGCM(block)

	 if err != nil {
		return nil, err
	 }

	 nonceSize := aesGCM.NonceSize()
	 nonce, cipherTextBytes := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]
	 return aesGCM.Open(nil, nonce, cipherTextBytes, nil)
}