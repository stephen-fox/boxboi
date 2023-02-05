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
	"os/signal"
	"path/filepath"

	"gitlab.com/stephen-fox/boxboi/internal/osspecific"
)

func main() {
	log.SetFlags(0)

	err := mainWithError()
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func mainWithError() error {
	ctx, cancelFn := signal.NotifyContext(context.Background(),
		osspecific.QuitSignals()...)
	defer cancelFn()

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get exe path - %s", err)
	}

	err = os.Chdir(filepath.Dir(exePath))
	if err != nil {
		return fmt.Errorf("failed to chdir to exe parent dir - %w",
			err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:3249")
	if err != nil {
		return err
	}
	defer listener.Close()

	errs := make(chan error, 1)
	go func() {
		var cancelCurrentClientFn func()
		defer func() {
			if cancelCurrentClientFn != nil {
				cancelCurrentClientFn()
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				errs <- err
				return
			}

			if cancelCurrentClientFn != nil {
				cancelCurrentClientFn()
			}

			var clientCtx context.Context
			clientCtx, cancelCurrentClientFn = context.WithCancel(ctx)

			go handleClient(clientCtx, conn)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errs:
		return err
	}
}

func handleClient(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	const help = `available commands:

> exit  - disconnect
> help  - show this information
> on    - boot boxboi
> off   - turn boxboi off`

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
