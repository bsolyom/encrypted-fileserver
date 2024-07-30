package handlers

import (
	"log"
	"os"
	"strconv"
	"time"
)

func isFileExpired(fileName string, expireAfter int) bool {
	stats, err := os.Stat("./storage/" + fileName)
	if err != nil {
		log.Println(err)
		return false
	}

	expireTime := stats.ModTime().AddDate(0, 0, expireAfter)
	return time.Now().After(expireTime)
}

func Cleanup() {
	for {
		var isExecuted bool

		log.Println("Searching for expired files...")

		expireAfter, _ := strconv.Atoi(os.Getenv("EXPIRE"))

		filesDir, err := os.Open("./storage")
		if err != nil {
			log.Println(err)
			return
		}
		defer filesDir.Close()

		files, err := filesDir.ReadDir(0)
		if err != nil {
			log.Println(err)
			return
		}

		for _, file := range files {
			if isFileExpired(file.Name(), expireAfter) {
				isExecuted = true
				os.Remove("./storage/" + file.Name())
			}
		}

		if isExecuted {
			log.Println("Expired files deleted")
		} else {
			log.Println("No expired files found")
		}
		time.Sleep(24 * time.Hour) // Checks for expired files every 24 hour
	}
}
