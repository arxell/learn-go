package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"log"
	"net/http"
	"strconv"
	"time"
)

var DB = make(map[string]string)

type ApiAuthKey struct {
	ID        uint
	Key       string
	Is_Active bool
	Rps       int
}

func (ApiAuthKey) TableName() string {
	return "partners_apiauthkey"
}

type Config struct {
	DatabaseHost     string `env:"DatabaseHost"`
	DatabasePort     int    `env:"PORT" envDefault:"5432"`
	DatabaseName     string `env:"DatabaseName" envDefault:"partner_20160602"`
	DatabaseUser     string `env:"DatabaseUser" envDefault:"postgres"`
	DatabasePassword string `env:"DatabasePassword"`
}

func main() {

	cfg := Config{}
	env.Parse(&cfg)
	fmt.Println(cfg)

	DB, err := gorm.Open(
		"postgres",
		"host="+cfg.DatabaseHost+
			" port="+strconv.Itoa(cfg.DatabasePort)+
			" user="+cfg.DatabaseUser+
			" dbname="+cfg.DatabaseName+
			" password="+cfg.DatabasePassword)

	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		ois := []ApiAuthKey{}
		DB.Limit(10).Find(&ois)
		fmt.Println(ois)
		c.String(200, "pong")
	})

	// Initalize
	Admin := admin.New(&qor.Config{DB: DB})

	// Create resources from GORM-backend model
	Admin.AddResource(&ApiAuthKey{})

	// Binding qor admin with Gin
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)
	router.Any("/admin/*w", gin.WrapH(mux))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
