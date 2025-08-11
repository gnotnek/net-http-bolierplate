package api

import (
	"context"
	"fmt"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/auth"
	"net-http-boilerplate/internal/category"
	"net-http-boilerplate/internal/config"
	"net-http-boilerplate/internal/pkg/encrypt"
	"net-http-boilerplate/internal/pkg/jwt"
	"net-http-boilerplate/internal/pkg/postgres"
	"net-http-boilerplate/internal/pkg/validator"
	"net-http-boilerplate/internal/post"
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

	// encrypt
	err := encrypt.Init(cfg.AppConfig.AppSalt, cfg.AppConfig.AppSaltIV, cfg.AppConfig.AppEncryptMethod)
	if err != nil {
		log.Panic().Err(err).Msgf("something went wrong: %s", err)
		panic(err)
	}

	// validator
	validator := validator.NewValidator()

	// database
	db := postgres.NewGORM(&cfg.Database)
	postgres.Migrate(db)

	// Initialize JWT service
	jwtService := jwt.NewJWT(cfg.JWT)

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(jwtService)

	// Repo
	userRepo := user.NewUserRepository(db)
	postRepo := post.NewPostRepository(db)
	categoryRepo := category.NewCategoryRepository(db)

	// Service
	userService := user.NewUserService(userRepo, jwtService)
	postService := post.NewPostService(postRepo)
	categoryService := category.NewCategoryService(categoryRepo)

	// Handler
	userHandler := user.NewUserHandler(userService, validator)
	postHandler := post.NewPostHandler(postService)
	categoryHandler := category.NewCategoryHandler(categoryService)

	r := chi.NewRouter()

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		resp.WriteJSON(w, http.StatusOK, "Pong")
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthRequired)

		// Posts
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", postHandler.FindAll)
			r.Post("/", postHandler.Create)
			r.Get("/{id}", postHandler.FindByID)
			r.Put("/{id}", postHandler.Update)
			r.Delete("/{id}", postHandler.Delete)
		})

		// Categories
		r.Route("/category", func(r chi.Router) {
			r.Get("/", categoryHandler.GetCategories)
			r.Get("/{id}", categoryHandler.GetCategory)
			r.Post("/", categoryHandler.CreateCategory)
			r.Put("/{id}", categoryHandler.UpdateCategory)
			r.Delete("/{id}", categoryHandler.DeleteCategory)
		})
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
