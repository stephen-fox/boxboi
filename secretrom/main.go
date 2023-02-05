package main

import (
	"bytes"
	"crypto/rc4"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	//go:embed key_hex
	keyHexStr string

	//go:embed last_four_bytes_hex
	expLastFourBytesHexStr string
)

func main() {
	log.SetFlags(0)

	err := mainWithError()
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func mainWithError() error {
	key, err := hex.DecodeString(keyHexStr)
	if err != nil {
		return fmt.Errorf("failed to hex-decode rc4 key - %w", err)
	}

	rc4Cipher, err := rc4.NewCipher(key)
	if err != nil {
		return err
	}

	expLastFourBytes, err := hex.DecodeString(expLastFourBytesHexStr)
	if err != nil {
		return fmt.Errorf("failed to hex-decode last four bytes - %w",
			err)
	}

	flashCiphertext, err := os.ReadFile("./flash")
	if err != nil {
		return err
	}

	flashLen := len(flashCiphertext)
	if flashLen < 4 {
		return errors.New("flash is less than four bytes")
	}

	flashPlaintext := make([]byte, flashLen)

	rc4Cipher.XORKeyStream(flashPlaintext, flashCiphertext)

	if !bytes.Equal(flashPlaintext[flashLen-4:], expLastFourBytes) {
		return errors.New("last four bytes of flash are incorrect")
	}

	kernelExe, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create kernel tmp file - %w", err)
	}
	defer kernelExe.Close()

	_, err = kernelExe.Write(flashPlaintext)
	if err != nil {
		return fmt.Errorf("failed to write kern to tmp file - %w", err)
	}

	err = kernelExe.Chmod(0700)
	if err != nil {
		return fmt.Errorf("failed to chmod kernel - %w", err)
	}

	_, err = os.Stdout.WriteString(kernelExe.Name())
	if err != nil {
		return fmt.Errorf("failed to write kernel path to stdout - %w",
			err)
	}

	return nil
}
