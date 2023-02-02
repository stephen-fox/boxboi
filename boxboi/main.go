// The entrypoint for boxboi :)
package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func main() {
	log.SetFlags(0)

	err := mainWithError()
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func mainWithError() error {
	err := resetChallenge()
	if err != nil {
		return fmt.Errorf("failed to reset challenge - %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:3249")
	if err != nil {
		return err
	}
	defer listener.Close()

	var cancelCurrentClientFn func()
	defer func() {
		if cancelCurrentClientFn != nil {
			cancelCurrentClientFn()
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		if cancelCurrentClientFn != nil {
			cancelCurrentClientFn()
		}

		var clientCtx context.Context
		clientCtx, cancelCurrentClientFn = context.WithCancel(
			context.Background())

		go handleClient(clientCtx, conn)
	}
}

func handleClient(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	const help = `available commands:

> exit  - disconnect
> help  - show this information
> on    - boot boxboi
> off   - turn boxboi off
> reset - reset the challenge to its original state (if something got broke)`

	// write is a helper function for writing messages to
	// a connected user.
	write := func(msg string) error {
		_, err := conn.Write([]byte(msg + "\n"))
		return err
	}

	var err error
	var kernel *exec.Cmd
	defer func() {
		if kernel != nil {
			_ = kernel.Process.Kill()
		}
	}()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()

		switch line {
		case "exit":
			return nil
		case "help":
			err = write(help)
			if err != nil {
				return err
			}
		case "on":
			if kernel != nil {
				err = write("already powered on")
				if err != nil {
					return err
				}

				continue
			}

			err = write("powering on...")
			if err != nil {
				return err
			}

			var poErr error

			kernel, poErr = powerOn(ctx, conn)
			if poErr != nil {
				err = write(poErr.Error())
				if err != nil {
					return err
				}
			}
		case "off":
			if kernel != nil {
				_ = kernel.Process.Kill()
				kernel = nil
			}

			err = write("powered off")
			if err != nil {
				return err
			}
		case "reset":
			err := resetChallenge()
			if err != nil {
				_ = write("failed to reset challenge - " +
					err.Error())
				return err
			} else {
				_ = write("challenge has been reset - " +
					"please reconnect")
				return nil
			}
		default:
			err = write("unknown command")
			if err != nil {
				return err
			}
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}

func powerOn(ctx context.Context, conn io.ReadWriter) (*exec.Cmd, error) {
	sr := exec.CommandContext(ctx, "./secretrom")
	stdout := bytes.NewBuffer(nil)
	sr.Stdout = stdout
	stderr := bytes.NewBuffer(nil)
	sr.Stderr = stderr

	err := sr.Run()
	if err != nil {
		return nil, fmt.Errorf("secretrom failure - %s - %w",
			stderr, err)
	}

	kernel := exec.CommandContext(ctx, stdout.String())
	kernel.Stderr = conn
	kernel.Stdout = conn

	err = kernel.Start()
	if err != nil {
		return nil, fmt.Errorf("kernel failure - %w", err)
	}

	return kernel, nil
}

func resetChallenge() error {
	backup, err := os.ReadFile("./flash.backup")
	if err != nil {
		return fmt.Errorf("failed to open flash backup - %w", err)
	}

	err = os.WriteFile("./flash", backup, 0o666)
	if err != nil {
		return fmt.Errorf("failed to overwrite flash with backup - %w",
			err)
	}

	err = os.Chmod("./flash", 0o666)
	if err != nil {
		return fmt.Errorf("failed to chmod flash - %w", err)
	}

	return nil
}
