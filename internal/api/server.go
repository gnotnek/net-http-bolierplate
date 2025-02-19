package api

import (
	"context"
	"fmt"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/auth"
	"net-http-boilerplate/internal/config"
	"net-http-boilerplate/internal/jwt"
	"net-http-boilerplate/internal/post"
	"net-http-boilerplate/internal/postgres"
	"net-http-boilerplate/internal/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func NewServer() *Server {
	cfg := config.Load()
	db := postgres.NewGORM(&cfg.Database)
	postgres.Migrate(db)

	// Initialize JWT service
	jwtService := jwt.NewJWT(cfg.JWT)

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(jwtService)

	// Repo
	userRepo := user.NewUserRepository(db)
	postRepo := post.NewPostRepository(db)

	// Service
	userService := user.NewUserService(userRepo)
	postService := post.NewPostService(postRepo)

	// Handler
	userHandler := user.NewUserHandler(userService, jwtService)
	postHandler := post.NewPostHandler(postService)

	r := chi.NewRouter()

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		resp.WriteJSON(w, http.StatusOK, "Pong")
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.RefreshToken)
	})

	r.With(authMiddleware.AuthRequired).Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.FindAll)
		r.Post("/", postHandler.Create)
		r.Get("/{id}", postHandler.FindByID)
		r.Put("/{id}", postHandler.Update)
		r.Delete("/{id}", postHandler.Delete)
	})

	return &Server{
		router: r,
	}
}

type Server struct {
	router *chi.Mux
}

// Run method of the Server struct runs the HTTP server on the specified port. It initializes
// a new HTTP server instance with the specified port and the server's router.
func (s *Server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	h := chainMiddleware(
		s.router,
		recoverHandler,
		loggerHandler(func(w http.ResponseWriter, r *http.Request) bool { return r.URL.Path == "/" }),
		realIPHandler,
		requestIDHandler,
		corsHandler,
	)

	httpServer := http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("Server forced to shutdown")
		}

		close(done)
	}()

	log.Info().Msgf("server serving on port %d", port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}

	<-done
	log.Info().Msg("Server stopped")
}
