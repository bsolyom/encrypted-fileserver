package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bsolyom/encrypted-fileserver/handlers"
	"github.com/joho/godotenv"
)

func setFileHeaders(w http.ResponseWriter, fileName string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
}

func handleDownloadWithDecrypt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 method not supported", http.StatusMethodNotAllowed)
		return
	}

	code := r.FormValue("code")
	psw := r.FormValue("psw")
	if psw == "" {
		http.Redirect(w, r, r.URL.Host+"/down/"+code, http.StatusSeeOther)
		return
	}

	filePath, err := handlers.DownloadFile(w, r, code, psw)
	if err != nil {
		log.Println(err)
		return
	}

	fileName := filepath.Base(filePath)
	setFileHeaders(w, fileName)
	log.Println("Encrypted file downloaded:", fileName)

	defer os.Remove(filePath) // Remove from the temp dir

	http.ServeFile(w, r, filePath)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 method not supported", http.StatusMethodNotAllowed)
		return
	}

	filePath, err := handlers.DownloadFile(w, r, filepath.Base(r.URL.Path), "")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := filepath.Base(filePath)
	setFileHeaders(w, fileName)
	log.Println("File downloaded:", fileName)

	http.ServeFile(w, r, filePath)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code, err := handlers.ReceiveFile(r, w)
	if err != nil {
		log.Println(err)
		if errors.As(err, new(*http.MaxBytesError)) {
			http.Error(w, "413 upload failed: the file has reached the maximum upload limit", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}
		return
	} else {
		log.Println("File uploaded:", code)
		fmt.Fprintln(w, code)
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load config file", err)
	}

	go handlers.Cleanup()

	portStr := ":" + os.Getenv("PORT")
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/down/", handleDownload)
	http.HandleFunc("/down", handleDownloadWithDecrypt)
	http.HandleFunc("/up", handleUpload)
	fmt.Printf("Starting server (%s) \n", portStr)
	if os.Getenv("CERTFILE") != "" && os.Getenv("KEYFILE") != "" {
		err := http.ListenAndServeTLS(portStr, os.Getenv("CERTFILE"), os.Getenv("KEYFILE"), nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := http.ListenAndServe(portStr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
