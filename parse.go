package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Info struct {
	name string
	url  string
}

// var url []Info
func Parse() {

	url := []Info{
		{
			name: "skins.json",
			url:  "https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/skins.json",
		},
		{
			name: "stickers.json",
			url:  "https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/stickers.json",
		},
		{
			name: "keychains.json",
			url:  "https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/keychains.json",
		},
		{
			name: "collections.json",
			url:  "https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/collections.json",
		},
		{
			name: "crates.json",
			url:  "https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/crates.json",
		},
	}
	for _, v := range url {
		time.Sleep(1 * time.Second)
		resp, err := http.Get(v.url)
		if err != nil {
			log.Printf("Error getting %s from %s: %v", v.name, v.url, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Printf("HTTP error %d for %s", resp.StatusCode, v.name)
			continue
		}

		filePath := filepath.Join("data", v.name)
		file, err := os.Create(filePath)
		if err != nil {
			resp.Body.Close()
			log.Printf("Failed to create file %s: %v", filePath, err)
			continue
		}

		size, err := io.Copy(file, resp.Body)

		file.Close()
		resp.Body.Close()

		if err != nil {
			log.Printf("Failed to save file %s: %v", filePath, err)
			continue
		}
		fmt.Printf("%s saved successfully (%d bytes)\n", v.name, size)
	}

}
