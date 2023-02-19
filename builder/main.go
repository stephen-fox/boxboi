// Tool that compiles boxboi's various components.
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	log.SetFlags(0)

	err := mainWithError()
	if err != nil {
		log.Fatalln(err)
	}
}

func mainWithError() error {
	ctx, cancelFn := signal.NotifyContext(context.Background(),
		os.Interrupt)
	defer cancelFn()

	err := os.MkdirAll("build", 0755)
	if err != nil {
		return err
	}

	err = goBuild(ctx, "kernel/main.go", "kernel/kernel")
	if err != nil {
		return err
	}

	kernel, err := os.ReadFile("kernel/kernel")
	if err != nil {
		return err
	}

	lastFourBytes := kernel[len(kernel)-4:]

	key := make([]byte, 256)
	_, err = rand.Read(key)
	if err != nil {
		return err
	}

	keyHexStr := hex.EncodeToString(key)

	err = os.WriteFile("secretrom/key_hex", []byte(keyHexStr), 0600)
	if err != nil {
		return err
	}

	lastFourBytesHexStr := hex.EncodeToString(lastFourBytes)

	err = os.WriteFile(
		"secretrom/last_four_bytes_hex",
		[]byte(lastFourBytesHexStr),
		0600)
	if err != nil {
		return err
	}

	encryptor := exec.CommandContext(ctx, "go", "run", "encrypt/main.go")
	encryptor.Env = append(os.Environ(), "RC4_KEY_HEX="+keyHexStr)
	encryptor.Stdin = bytes.NewReader(kernel)
	encryptor.Stderr = os.Stderr
	ct := bytes.NewBuffer(nil)
	encryptor.Stdout = ct

	err = encryptor.Run()
	if err != nil {
		return err
	}

	err = os.WriteFile("build/flash", ct.Bytes(), 0666)
	if err != nil {
		return err
	}

	err = os.Chmod("build/flash", 0666)
	if err != nil {
		return err
	}

	err = goBuild(ctx, "secretrom/main.go", "build/secretrom")
	if err != nil {
		return err
	}

	err = goBuild(ctx, "boxboi/main.go", "build/boxboi")
	if err != nil {
		return err
	}

	return nil
}

func goBuild(ctx context.Context, filePath string, outputPath string) error {
	goBuild := exec.CommandContext(ctx,
		"go", "build", "-o", outputPath, filePath)
	goBuild.Stderr = os.Stderr
	goBuild.Stdout = os.Stdout
	goBuild.Env = os.Environ()

	if targetOS := os.Getenv("_GOOS"); targetOS != "" {
		goBuild.Env = append(goBuild.Env, "GOOS="+targetOS)
		log.Printf("[build %s]: GOOS set to '%s'",
			filePath, targetOS)
	}

	if targetArch := os.Getenv("_GOARCH"); targetArch != "" {
		goBuild.Env = append(goBuild.Env, "GOARCH="+targetArch)
		log.Printf("[build %s]: GOARCH set to '%s'",
			filePath, targetArch)
	}

	log.Printf("[build %s]: exec: '%s'", filePath, goBuild.String())

	return goBuild.Run()
}
