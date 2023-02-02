package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(0)

	err := mainWithError()
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func mainWithError() error {
	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFn()

	// https://textkool.com/en/ascii-art-generator
	_, err := os.Stdout.WriteString(`___.                   ___.            .__
\_ |__    ____ ___  ___\_ |__    ____  |__|
 | __ \  /  _ \\  \/  / | __ \  /  _ \ |  |
 | \_\ \(  <_> )>    <  | \_\ \(  <_> )|  |
 |___  / \____//__/\_ \ |___  / \____/ |__|
     \/              \/     \/

(please insert a game disc)
`)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
