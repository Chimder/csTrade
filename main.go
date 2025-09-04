package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func main() {
	message := []byte("Hello!")
	secret := []byte("superkey")

	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	signature := mac.Sum(nil)


	fmt.Printf("HMAC signature: %x\n", signature)
}
