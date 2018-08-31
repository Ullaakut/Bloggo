package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ullaakut/Bloggo/controller"
	"github.com/Ullaakut/Bloggo/logger"
	"github.com/Ullaakut/Bloggo/repo"
	"github.com/Ullaakut/Bloggo/service"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {

	log := logger.NewZeroLog(os.Stderr)
	log.Info().Msg("bloggo is barking up")

	config := GetConfig()
	config.Print(log)

	if config.jwtSecret == "" {
		log.Fatal().Msg("JWT secret not set. please set it in the environment for bloggo to work properly")
		os.Exit(1)
	}

	zerolog.SetGlobalLevel(logger.ParseLevel(config.LogLevel))

	// Catch signals
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// Use zerolog for debugging HTTP requests
	e.Logger.SetLevel(5) // Disable default logging
	e.Use(logger.HTTPLogger(log))

	// Initialize the database
	// Retry until it is successful or the retryDuration is over
	var db *gorm.DB
	var err error
	startTime := time.Now()
	try(log, config.MySQLRetryInterval, func() error {
		db, err = gorm.Open("mysql", config.MySQLURL)
		return err
	}, func() bool {
		return time.Now().Sub(startTime) < config.MySQLRetryDuration
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could not initialize mysql connection")
		os.Exit(1)
	}
	log.Info().Msg("connection to mysql successful")

	blogPostRepository := repo.NewBlogPostRepositoryMySQL(log, db)
	userRepository := repo.NewUserRepositoryMySQL(log, db)

	accessService := service.NewAccess(log, userRepository, config.jwtSecret)
	tokenService := service.NewToken(log, userRepository, config.jwtSecret)

	blogController := controller.NewBlog(log, blogPostRepository)
	userController := controller.NewUser(log, userRepository, tokenService)
	authController := controller.NewAuth(log, accessService)

	// Bind routes to controller methods

	// Login&Registration API
	e.POST("/register", userController.Register)
	e.POST("/login", userController.Login)

	// Blog post API
	e.POST("/posts", blogController.Create, authController.Authorize)
	e.GET("/posts", blogController.ReadAll)
	e.GET("/posts/:id", blogController.Read)
	e.PUT("/posts/:id", blogController.Update, authController.Authorize)
	e.DELETE("/posts/:id", blogController.Delete, authController.Authorize)

	// Graceful enables graceful shutdown of the HTTP server
	e.Server.Addr = fmt.Sprintf("%v:%v", config.ServerAddress, config.ServerPort)
	server := &graceful.Server{
		NoSignalHandling: true,
		Server:           e.Server,
	}

	// Start server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Info().Err(err).Msg("could not start server")
			os.Exit(1)
		}
	}()

	log.Info().Msg("bloggo is up")

	// Wait for server to be stopped
	<-sig
	signal.Stop(sig)
	close(sig)

	log.Info().Msg("bloggo is shutting down")

	server.Stop(15 * time.Second)

	log.Info().Msg("bloggo shutdown complete")

	os.Exit(0)
}

// try tries to execute a given function
// if it fails, it will keep retrying until the given shouldRetry function returns false
func try(logger *zerolog.Logger, retryDelay time.Duration, fn func() error, shouldRetry func() bool) error {
	for {
		err := fn()
		if err == nil {
			return nil
		}

		if !shouldRetry() {
			logger.Error().Err(err).Msg("operation failed too many times, aborting")
			return err
		}

		logger.Error().Err(err).Msg("operation failed, will retry")
		time.Sleep(retryDelay)
	}
}
