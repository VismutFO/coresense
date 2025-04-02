package api

/*
import (
	"context"
	"time"

	"gorm.io/gorm"
	"net/http"
	"strconv"

	"coresense/app/server/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

var jwtSecret = []byte("your_secret_key")

type Processor struct {
	config config.Server
	db     *gorm.DB
	logger zerolog.Logger
	mux    *http.ServeMux
}

// NewProcessor initializes the processor.
func NewProcessor(ctx context.Context, config config.Server, db *gorm.DB, logger zerolog.Logger) (*Processor, error) {
	logger.Debug().Ctx(ctx).Msg("In")
	mux := http.NewServeMux()
	processor := &Processor{
		config: config,
		db:     db,
		logger: logger,
		mux:    mux,
	}

	processor.routes(ctx)
	logger.Debug().Ctx(ctx).Msg("Out")
	return processor, nil
}

// routes sets up API endpoints.
func (p *Processor) routes(ctx context.Context) {
	p.logger.Debug().Ctx(ctx).Msg("In")
	p.mux.Handle("/register", http.HandlerFunc(p.RegisterUser))
	p.mux.Handle("/login", http.HandlerFunc(p.AuthorizeUser))
	p.mux.Handle("/registerCustomer", http.HandlerFunc(p.RegisterBusinessCustomer))
	p.mux.Handle("/loginCustomer", http.HandlerFunc(p.AuthorizeBusinessCustomer))
	p.mux.Handle("/addServiceTemplate", http.HandlerFunc(p.AddServiceTemplate))
	p.mux.Handle("/addFilledService", http.HandlerFunc(p.AddFilledService))
	p.mux.Handle("/checkFilledServices", http.HandlerFunc(p.GetFilledServices))

	p.logger.Debug().Ctx(ctx).Msg("Out")
}

// generateJWT creates a new JWT token for a user.
func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// Start runs the HTTP server.
func (p *Processor) Start(ctx context.Context) error {
	p.logger.Debug().Ctx(ctx).Msg("In")
	return http.ListenAndServe(":"+strconv.Itoa(p.config.HTTP.Port), p.mux)
}
*/
