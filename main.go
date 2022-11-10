package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AdiKhoironHasan/mister-aladin-articles/pkg/database"

	integ "github.com/AdiKhoironHasan/mister-aladin-articles/internal/integration"
	SqlRepo "github.com/AdiKhoironHasan/mister-aladin-articles/internal/repository/postgresql"
	NoSqlRepo "github.com/AdiKhoironHasan/mister-aladin-articles/internal/repository/redis"
	"github.com/AdiKhoironHasan/mister-aladin-articles/internal/services"
	handlers "github.com/AdiKhoironHasan/mister-aladin-articles/internal/transport/http"
	"github.com/AdiKhoironHasan/mister-aladin-articles/internal/transport/http/middleware"

	"github.com/apex/log"
	"github.com/labstack/echo"

	"github.com/spf13/viper"
)

func main() {

	errChan := make(chan error)

	e := echo.New()
	m := middleware.NewMidleware()

	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config-dev")

	err := viper.ReadInConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	dbhost, dbUser, dbPassword, dbName, dbPort :=
		viper.GetString("db.postgre.host"),
		viper.GetString("db.postgre.user"),
		viper.GetString("db.postgre.password"),
		viper.GetString("db.postgre.dbname"),
		viper.GetString("db.postgre.port")

	PostgreSqlDB, err := database.PostgreSqllInitialize(dbhost, dbUser, dbPassword, dbName, dbPort)

	if err != nil {
		log.Fatal("Failed to Connect Postgre SQL Database: " + err.Error())
	}

	dbhost, dbUser, dbPassword, dbPort =
		viper.GetString("db.redis.host"),
		viper.GetString("db.redis.user"),
		viper.GetString("db.redis.password"),
		viper.GetString("db.redis.port")

	RedisDB, err := database.RedislInitialize(dbhost, dbUser, dbPassword, dbPort)
	if err != nil {
		log.Fatal("Failed to Connect Redis Database: " + err.Error())
	}

	defer func() {
		err := PostgreSqlDB.Conn.Close()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = RedisDB.Conn.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	e.Use(m.CORS)

	sqlrepo := SqlRepo.NewRepo(PostgreSqlDB.Conn)
	noSqlRepo := NoSqlRepo.NewRepo(RedisDB.Conn)
	integSrv := integ.NewService()
	srv := services.NewService(sqlrepo, noSqlRepo, integSrv)
	// srv := services.NewService(sqlrepo)
	handlers.NewHttpHandler(e, srv)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errChan <- e.Start(":" + viper.GetString("server.port"))
	}()

	e.Logger.Print("Starting ", viper.GetString("appName"))
	err = <-errChan
	log.Error(err.Error())

}
