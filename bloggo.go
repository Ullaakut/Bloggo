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

	// TODO: Config override via config file, env, consul kv...
	config := DefaultConfig()
	config.JWTSecret = os.Getenv("BLOGGO_JWT_SECRET")
	config.Print(log)

	if config.JWTSecret == "" {
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

	db, err := connectMySQL(log, config.MySQLURL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not initialize mysql connection")
		os.Exit(1)
	}

	blogPostRepository := repo.NewBlogPostRepositoryMySQL(log, db)
	userRepository := repo.NewUserRepositoryMySQL(log, db)

	accessService := service.NewAccess(log, userRepository, config.TrustedSource, config.JWTSecret)

	blogController := controller.NewBlog(log, blogPostRepository)
	authController := controller.NewAuth(log, accessService)

	// Bind route to controller method
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

	server.Stop(config.GracefulShutdownTimeout)

	log.Info().Msg("bloggo shutdown complete")

	os.Exit(0)
}

func connectMySQL(log *zerolog.Logger, url string) (*gorm.DB, error) {
	// Setup DB connector
	// Try 50 times, in case the db is slow to start
	connectionAttempts := 0
	var db *gorm.DB
	for {
		var err error
		db, err = gorm.Open("mysql", url)
		if err == nil {
			break
		}

		if connectionAttempts > 50 {
			return db, err
		}
		connectionAttempts++

		log.Debug().Msg("failed to connect to mysql database, will retry in 2 seconds...")
		<-time.After(2 * time.Second)
	}

	log.Debug().Msg("mysql connection successful")
	return db, nil
}
