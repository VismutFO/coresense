package clientapi

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	metricsprometheus "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsmiddleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"coresense/pkg/core/model/constants"
)

type zerologWrapper struct {
	logger zerolog.Logger
}

func (zw zerologWrapper) Write(p []byte) (n int, err error) {
	zw.logger.Error().Msg(string(p))
	return len(p), nil
}

func Start(ctx context.Context, username, password, certificate, privateKey, netInterface string, port int) error {
	logger := zerolog.Nop() // stub

	service := web.NewService(openapi31.NewReflector())

	service.OpenAPISchema().SetTitle("Client API")
	service.OpenAPISchema().SetDescription("This service provides API for end users.")
	service.OpenAPISchema().SetVersion("v1.0.0")

	service.Use(
		middleware.Recoverer,
		std.HandlerProvider("", metricsmiddleware.New(metricsmiddleware.Config{
			Recorder: metricsprometheus.NewRecorder(metricsprometheus.Config{}),
		})),
		middleware.Heartbeat("/check"),
		cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodPost,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
		}).Handler,
	)

	service.Handle("/", http.NotFoundHandler())

	// Prepare middleware with suitable security schema.
	// It will perform actual security check for every relevant request.
	apiAuth := middleware.BasicAuth("Restricted Access", map[string]string{username: password})

	// Prepare API schema updater middleware.
	// It will annotate handler documentation with security schema.
	apiSecuritySchema := nethttp.HTTPBasicSecurityMiddleware(service.OpenAPICollector, "API", "API access")

	service.Route("/api", func(r chi.Router) {
		r.With(apiAuth, apiSecuritySchema).Route("/v{version}", func(r chi.Router) {
			r.Method(http.MethodPost, "/repost", nethttp.NewHandler(func() usecase.Interactor {
				u := usecase.NewInteractor(func(ctx context.Context, request struct {
					// BatchID              string         `json:"batch_id"     form:"batch_id"`
					// OperationID          string         `json:"operation_id" form:"operation_id"`
					Version              string         `json:"-"            example:"1.0.0"     path:"version"`
					AdditionalProperties map[string]any `json:"-"` // All unmatched properties.
				}, output *struct{},
				) error {
					/* batchID, err := uuid.Parse(request.BatchID)
					if err != nil || batchID == uuid.Nil {
						return status.Wrap(errors.New("missing or invalid batch id"), status.InvalidArgument)
					}

					operationID, err := uuid.Parse(request.OperationID)
					if err != nil || operationID == uuid.Nil {
						return status.Wrap(errors.New("missing or invalid operation id"), status.InvalidArgument)
					} */
					var err error
					if err != nil {
						var appErr interface {
							rest.ErrWithAppCode
							rest.ErrWithCanonicalStatus
						}

						if !errors.As(err, &appErr) {
							err = usecase.Error{
								StatusCode: status.Internal,
								Value:      err,
							}
						}

						logger.Err(err).Dict(constants.Args, zerolog.Dict().Interface("error", err)).Msg("usecase error")

						return err
					}

					return nil
				})
				u.SetTags("repost")
				u.SetTitle("Repost")
				u.SetDescription("Repost expects a batch UUID to restart batch processing.")
				u.SetExpectedErrors(status.InvalidArgument, status.Unavailable)

				return u
			}(), nethttp.SuccessStatus(http.StatusAccepted)))
		})
	})

	// Swagger UI should be provided by swgui handler constructor
	service.Docs("/docs", swgui.New)

	// start server
	// tlsConfig, err := config.LoadTLSConfig(conf)
	// if err != nil {
	//	return errors.New("failed to load TLS configuration")
	// }

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", netInterface, port),
		Handler: service,
		// TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ErrorLog:     log.New(zerologWrapper{logger: logger.With().Logger()}, "", 0),
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServeTLS(certificate, privateKey)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return nil
	}

	return nil
}
