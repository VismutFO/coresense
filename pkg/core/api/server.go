package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your_secret_key")

type Server struct {
	mux        *http.ServeMux // Shared router
	server     *http.Server
	processors []Processor // List of processors implementing business logic
}

type Processor interface {
	Routes(mux *http.ServeMux) // Method to add routes to the shared mux
}

// NewServer initializes the server and prepares the shared router
func NewServer() *Server {
	return &Server{
		mux:        http.NewServeMux(), // Shared router
		processors: []Processor{},      // Empty processor list
	}
}

// RegisterProcessor adds a new processor to the server
func (s *Server) RegisterProcessor(processor Processor) {
	s.processors = append(s.processors, processor)
	processor.Routes(s.mux) // Each processor registers its routes onto the shared mux
}

// Start runs the HTTP server with the shared router
func (s *Server) Start(port string) error {
	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: corsMiddleware(s.mux), // Apply CORS middleware to the shared router
	}

	return s.server.ListenAndServe() // Start the server
}

func (s *Server) Stop() {
	_ = s.server.Shutdown(context.Background())
}

// generateJWT creates a new JWT token for a user.
func generateJWT(username string, claims map[string]any) (string, error) {
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	return token.SignedString(jwtSecret)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Use "*" to allow all origins, or specify your Flutter app's domain
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle pre-flight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// authMiddleware checks the validity of the JWT token and ensures proper access
func authMiddleware(requiredRole string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is correct
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and validate them
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check for token expiration
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					http.Error(w, "Token has expired", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Token expiration not found", http.StatusUnauthorized)
				return
			}

			// Check for the proper role
			if role, ok := claims["role"].(string); ok {
				if role != requiredRole {
					http.Error(w, "Unauthorized role", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Role not found in token", http.StatusUnauthorized)
				return
			}

			// Add claims to the context
			ctx := context.WithValue(r.Context(), "claims", claims)

			// Token is valid and user has the correct role; proceed to the handler
			// Pass the modified request to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}
