package filecrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"os"
	"strings"
)

func DecryptFile(fileName string, psw string) (string, error) {
	key := GenerateKey(psw)

	ciphertext, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key[:16])
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	tempFile, err := os.Create(os.TempDir() + strings.Replace(fileName, "./storage", "", 1)) // os.CreateTemp isnt a good option because of the filename
	if err != nil {
		return "", err
	}

	_, err = tempFile.Write(ciphertext)
	if err != nil {
		tempFile.Close()
		return "", err
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		tempFile.Close()
		return "", err
	}

	return tempFile.Name(), nil
}
