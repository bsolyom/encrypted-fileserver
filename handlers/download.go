package handlers

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/bsolyom/encrypted-fileserver/filecrypto"
)

func loadCodes(filePath string) (map[string]string, error) {
	var codesMap = make(map[string]string)
	codesFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer codesFile.Close()

	scanner := bufio.NewScanner(codesFile)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		codesMap[parts[0]] = parts[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return codesMap, nil
}

func DownloadFile(w http.ResponseWriter, r *http.Request, code, key string) (string, error) {
	codesMap, err := loadCodes("codes.txt")
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return "", err
	}

	filePath, ok := codesMap[code]
	filePath = "./storage/" + filePath
	if !ok {
		http.Error(w, "404 file not found", http.StatusNotFound)
		return "", err
	}

	if key != "" {
		filePathEnc, err := filecrypto.DecryptFile(filePath, key)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return "", err
		}
		return filePathEnc, nil
	}
	return filePath, nil
}
