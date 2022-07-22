package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	bigtable "cloud.google.com/go/bigtable"
	deviceModel "github.com/device-auth/model"
	deviceDelivery "github.com/device-auth/implementation/delivery/http"
	deviceBigTableRepository "github.com/device-auth/implementation/repository/bigtable"
	devicePostgresRepository "github.com/device-auth/implementation/repository/postgresql"
	deviceUsecase "github.com/device-auth/implementation/usecase"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(`path: ` + path)
	viper.SetConfigType(`json`)
	viper.SetConfigName(`config`)
	viper.AddConfigPath(`./`)
	viper.AddConfigPath(`../`)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config file not found")
		}
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	fmt.Println("Go Time")

	flag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvPrefix(viper.GetString("ENV"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	timeoutContext := time.Duration(viper.GetInt("CONTEXT.TIMEOUT")) * time.Second

	url := viper.GetString("ENV_PPSA")
	if url == "" {
		fmt.Println("Configuration Error: ENV_PPSA address not available")
	}
	port := viper.GetString("ENV_OPA_PORT")
	if port == "" {
		fmt.Println("Configuration Error: ENV_OPA_PORT port not available")
	}
	dbType := viper.GetString("dbType")
	if dbType== "" {
		fmt.Println("Configuration Error: dbType not available")
	}
	var deviceRepository deviceModel.IDeviceRepository;
	if dbType=="psql"{
		log.Println("psql")
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
			log.Fatal(err)
		}
		err = dbConn.Ping()
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			err := dbConn.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		deviceRepository = devicePostgresRepository.NewDeviceRepository(dbConn)
	} else if dbType=="btable"{
		log.Println("btable")
		ctx := context.Background()
		bProject := viper.GetString("bProject")
		if bProject == "" {
			log.Fatalf("No bigTable Project Specified in Config")
		}
		bInstance := viper.GetString("bInstance")
		if bInstance == "" {
			log.Fatalf("No bigTable Instance Specified in Config")
		}
		client, err := bigtable.NewClient(ctx, bProject, bInstance)
		tbl := client.Open(tableName)
		deviceRepository = deviceBigTableRepository.NewDeviceRepository(client,table)
		if err != nil {
				log.Fatalf("Could not create data operations client: %v", err)
		}


	} else{
		log.Fatal("Db Type Not Found")
	}
	deviceUseCase := deviceUsecase.NewDeviceUsecase(deviceRepository, timeoutContext)
	deviceDelivery.NewDeviceHandler(e, deviceUseCase)

	log.Fatal(e.Start(viper.GetString("ENV_AUTH_SERVER")))
}
