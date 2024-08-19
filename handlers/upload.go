package handlers

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/bsolyom/encrypted-fileserver/filecrypto"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateCode(fileName string) (string, error) {
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	codesFile, err := os.OpenFile("codes.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer codesFile.Close()

	if _, err := codesFile.WriteString(string(code) + ":" + fileName + "\n"); err != nil {
		return "", err
	}

	return string(code), nil
}

func ReceiveFile(r *http.Request, w http.ResponseWriter) (string, error) {
	maxSize, _ := strconv.ParseInt(os.Getenv("SIZELIMIT"), 10, 64)

	r.Body = http.MaxBytesReader(w, r.Body, maxSize*1000000)
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	psw := r.FormValue("psw")

	if psw != "" {
		file, err = filecrypto.EncryptFile(file, psw)
		if err != nil {
			return "", err
		}
	}

	dst, err := os.Create("./storage/" + header.Filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	code, err := generateCode(header.Filename)
	if err != nil {
		return "", err
	}

	return code, nil
}
