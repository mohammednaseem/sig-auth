package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	bigtable "cloud.google.com/go/bigtable"
	deviceDelivery "github.com/device-auth/implementation/delivery/http"
	deviceBigTableRepository "github.com/device-auth/implementation/repository/bigtable"
	devicePostgresRepository "github.com/device-auth/implementation/repository/postgresql"
	deviceUsecase "github.com/device-auth/implementation/usecase"
	deviceModel "github.com/device-auth/model"

	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	Log "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	lecho "github.com/ziflex/lecho"
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	log.Error().Err(err).Msg((`path: ` + path))
	viper.SetConfigType(`json`)
	viper.SetConfigName(`config`)
	viper.AddConfigPath(`./`)
	viper.AddConfigPath(`../`)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Error().Err(err).Msg("config file not found")
		}
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Info().Msg("Service RUN on DEBUG mode")
	}
}

func main() {
	log.Info().Msg("Go Time")

	flag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvPrefix(viper.GetString("ENV"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	e := echo.New()
	logger := lecho.New(
		os.Stdout,
		lecho.WithLevel(Log.DEBUG),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)
	e.Logger = logger
	e.Use(lecho.Middleware(lecho.Config{
		Logger: logger}))
	//e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	timeoutContext := time.Duration(viper.GetInt("CONTEXT.TIMEOUT")) * time.Second

	url := viper.GetString("ENV_PPSA")
	if url == "" {
		log.Error().Msg("Configuration Error: ENV_PPSA address not available")
	}
	port := viper.GetString("ENV_OPA_PORT")
	if port == "" {
		log.Error().Msg("Configuration Error: ENV_OPA_PORT port not available")
	}
	dbType := viper.GetString("dbType")
	if dbType == "" {
		log.Error().Msg("Configuration Error: dbType not available")
	}
	var deviceRepository deviceModel.IDeviceRepository
	if dbType == "psql" {
		log.Print("psql")
		dbHost := viper.GetString("ENV_DBHOST")
		if dbHost == "" {
			dbHost = viper.GetString(`db_host`)
		}
		dbPort := viper.GetString("ENV_PORT")
		if dbPort == "" {
			dbPort = viper.GetString(`db_port`)
		}
		dbUser := viper.GetString("ENV_USER")
		if dbUser == "" {
			dbUser = viper.GetString(`db_user`)
		}
		dbPass := viper.GetString("ENV_PASS")
		if dbPass == "" {
			dbPass = viper.GetString(`db_pass`)
		}
		dbName := viper.GetString("DB_NAME")
		if dbName == "" {
			dbName = viper.GetString(`db_name`)
		}

		// postgresql
		dbPortInt, _ := strconv.Atoi(dbPort)
		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPortInt, dbUser, dbPass, dbName)

		dbConn, err := sql.Open("postgres", connectionString)

		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		err = dbConn.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		defer func() {
			err := dbConn.Close()
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
		}()

		deviceRepository = devicePostgresRepository.NewDeviceRepository(dbConn)
	} else if dbType == "bigtable" {
		ctx := context.Background()
		bProject := viper.GetString("bigProject")
		if bProject == "" {
			log.Fatal().Msg("No bigTable Project Specified in Config")
		}
		bInstance := viper.GetString("bigInstance")
		if bInstance == "" {
			log.Fatal().Msg("No bigTable Instance Specified in Config")
		}
		bTable := viper.GetString("bTable")
		if bTable == "" {
			log.Fatal().Msg("No bigTable Table Specified in Config")
		}
		client, err := bigtable.NewClient(ctx, bProject, bInstance)
		log.Print("bigtable")
		tbl := client.Open(bTable)
		deviceRepository = deviceBigTableRepository.NewDeviceRepository(client, tbl)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create data operations client ")
		}

	} else {
		log.Fatal().Msg("Db Type Not Found")
	}
	deviceUseCase := deviceUsecase.NewDeviceUsecase(deviceRepository, timeoutContext)
	deviceDelivery.NewDeviceHandler(e, deviceUseCase)

	log.Fatal().Err((e.Start(viper.GetString("ENV_AUTH_SERVER")))).Msg("")
}
