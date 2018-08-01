package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ullaakut/Bloggo/controller"
	"github.com/Ullaakut/Bloggo/repo"
	"github.com/Ullaakut/Bloggo/service"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {

	log := NewZeroLog(os.Stderr)
	log.Info().Msg("bloggo is barking up")

	// TODO: Config override via config file, env, consul kv...
	config := DefaultConfig()
	config.Print(log)

	zerolog.SetGlobalLevel(parseLevel(config.LogLevel))

	// Catch signals
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// Use zerolog for debugging HTTP requests
	e.Logger.SetLevel(5) // Disable default logging
	e.Use(HTTPLogger(log))

	// Setup DB connector
	// db, err := gorm.Open("mysql", config.MySQLURL)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("could not initialize mysql connection")
	// 	os.Exit(1)
	// }

	blogPostRepository := repo.NewBlogPostRepository() //(db)
	userRepository := repo.NewUserRepository()         //(db)

	accessService := service.NewAccess(userRepository, config.TrustedSource)

	blogController := controller.NewBlog(blogPostRepository)
	authController := controller.NewAuth(accessService)

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
