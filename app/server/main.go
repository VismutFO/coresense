package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/sync/errgroup"

	"coresense/app/server/config"
	"coresense/pkg/common/utils/logging"
	"coresense/pkg/core/api"
	"coresense/pkg/core/model/collection/postgres/connection"
)

const envPrefix = "SERVER"

func main() {
	os.Exit(execute())
}

func execute() int {
	logger := logging.NewDefault()
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	server := api.NewServer()

	go func() {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals

		server.Stop()
		cancel()
	}()

	appConfig, err := config.New(envPrefix, logger)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("failed to read app config")
		return 1
	}

	db, err := connection.GetConnection(appConfig.Database, logger)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("failed to create database connection")
		return 1
	}

	eg.Go(func() error {
		server.RegisterProcessor(api.NewUserProcessor(appConfig.Server, db, logger))
		server.RegisterProcessor(api.NewBusinessProcessor(appConfig.Server, db, logger))

		if err := server.Start(strconv.Itoa(appConfig.Server.HTTP.Port)); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				return err
			}
		}
		return nil
	})

	/*
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
	*/
	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}
