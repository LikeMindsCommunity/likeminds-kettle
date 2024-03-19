package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/nateshr/likeminds-authentication/environment"
	"github.com/nateshr/likeminds-authentication/logging"
)

var keyPhrase = mdHashing(environment.GoDotEnvVariable("SECRET_KEY"))

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:])
}

func createAESBlock(keyPhrase []byte) cipher.Block {
	aesBlock, err := aes.NewCipher([]byte(keyPhrase))

	if err != nil {
		logging.Fatal(err)
		panic(err)
	}

	return aesBlock
}

func createGCM(aesBlock cipher.Block) cipher.AEAD {
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		panic(err)
	}

	return gcm
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return data
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(value []byte) string {
	keyPhrase := mdHashing(environment.GoDotEnvVariable("SECRET_KEY"))

	aesBlock := createAESBlock([]byte(keyPhrase))
	gcm := createGCM(aesBlock)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return ""
	}

	cipherText := gcm.Seal(nonce, nonce, value, nil)
	return Encode(cipherText)
}

// Decrypt method is to extract back the encrypted text
func Decrypt(cipheredText string) []byte {
	cipherText := Decode(cipheredText)

	aesBlock := createAESBlock([]byte(keyPhrase))
	gcm := createGCM(aesBlock)

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, _ := gcm.Open(nil, nonce, cipherText, nil)

	return plainText
}
