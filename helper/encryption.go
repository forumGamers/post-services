package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func Encryption(data string) string {
	byteMsg := []byte(data)
	block, err := aes.NewCipher([]byte(os.Getenv("ENCRPYTION_KEY")))
	if err != nil {
		PanicIfError(fmt.Errorf("could not create new cipher: %v", err))
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		PanicIfError(fmt.Errorf("could not encrypt: %v", err))
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText)
}

func Decryption(data string) string {
	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		PanicIfError(fmt.Errorf("could not base64 decode: %v", err))
	}

	block, err := aes.NewCipher([]byte(os.Getenv("ENCRPYTION_KEY")))
	if err != nil {
		PanicIfError(fmt.Errorf("could not create new cipher: %v", err))
	}

	if len(cipherText) < aes.BlockSize {
		panic(InvalidChiper.Error()) //buat testing untuk tes kalau enkripsinya ga sesuai
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}