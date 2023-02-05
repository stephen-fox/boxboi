// Encrypts stdin data using a RC4 key supplied via an environment variable and
// writes the resulting ciphertext to stdout.
package main

import (
	"crypto/rc4"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	err := mainWithError()
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func mainWithError() error {
	keyHexStr := os.Getenv("RC4_KEY_HEX")
	if keyHexStr == "" {
		return errors.New("please set the 'RC4_KEY_HEX' env var")
	}

	key, err := hex.DecodeString(keyHexStr)
	if err != nil {
		return fmt.Errorf("failed to hex-decode rc4 key - %w", err)
	}

	rc4Cipher, err := rc4.NewCipher(key)
	if err != nil {
		return err
	}

	pt, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	ct := make([]byte, len(pt))

	rc4Cipher.XORKeyStream(ct, pt)

	_, err = os.Stdout.Write(ct)
	return err
}
