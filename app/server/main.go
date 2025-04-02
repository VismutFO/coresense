package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"coresense/pkg/core/service/clientapi"
)

func main() {
	os.Exit(execute())
}

func execute() int {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := clientapi.Start(
			ctx,
			"username",
			"password",
			"cert.crt",
			"cert.key",
			"",
			8443,
		); err != nil {
			return err
		}
		return nil
	})

	go func() {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals

		// todo: shut down tracer
		cancel()
	}()

	if err := eg.Wait(); err != nil {
		return 1
	}

	return 0
}
