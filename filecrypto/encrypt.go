package filecrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"mime/multipart"
	"os"
)

func GenerateKey(psw string) [32]byte {
	return sha256.Sum256([]byte(psw)) // could be MD5, but its safer
}

func EncryptFile(file multipart.File, psw string) (*os.File, error) {
	key := GenerateKey(psw)

	plaintext, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key[:16])
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	buffer := bytes.NewBuffer(ciphertext)

	tempFile, err := os.CreateTemp("", "encrypted-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	_, err = buffer.WriteTo(tempFile)
	if err != nil {
		tempFile.Close()
		return nil, err
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		tempFile.Close()
		return nil, err
	}

	return tempFile, nil
}
