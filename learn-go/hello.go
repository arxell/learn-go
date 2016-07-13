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
	ID        int
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

func get_apiauthkeys(DB *gorm.DB) map[string]string {
	apiauthkeys := []ApiAuthKey{}
	DB.Find(&apiauthkeys)

	m := map[string]string{}
	for _, apiauthkey := range apiauthkeys {
		m[strconv.Itoa(apiauthkey.ID)] = apiauthkey.Key
	}
	return m
}

// Binding from JSON
type HotelRatesEndpointIn struct {
	Adults   int      `json:"adults" binding:"required"`
	Checkin  string   `json:"checkin" binding:"required"`
	Checkout string   `json:"checkout" binding:"required"`
	Ids      []string `json:"Ids" binding:"required"`
}

func main() {

	// Load env config
	cfg := Config{}
	env.Parse(&cfg)
	fmt.Println(cfg)

	// Init database
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

	// Init HTTP BASIC AUTH
	authorized := router.Group("/api", gin.BasicAuth(get_apiauthkeys(DB)))

	// === Endpoints part START ===
	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	// test json result
	router.GET("/ping/json", func(c *gin.Context) {
		// You also can use a struct
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		c.JSON(http.StatusOK, msg)
	})
	// hotel/rates
	authorized.GET("/ping/auth", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		adults := c.DefaultQuery("adults", "2")
		checkin := c.Query("checkin")
		checkout := c.Query("checkout")

		fmt.Println(user)
		c.JSON(http.StatusOK, gin.H{
			"user":     user,
			"adults":   adults,
			"checkin":  checkin,
			"checkout": checkout,
		})
	})
	// hotel/rates
	authorized.POST("/affiliate/v2/hotel/rates", func(c *gin.Context) {
		var params HotelRatesEndpointIn
		if c.BindJSON(&params) == nil {
			c.JSON(http.StatusOK, gin.H{"params": params})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "error"})
		}
	}) // === Endpoints part END ===

	//	=== Admin part START===
	Admin := admin.New(&qor.Config{DB: DB})
	//	Create resources from GORM-backend model
	Admin.AddResource(&ApiAuthKey{})
	//	Binding admin with Gin
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)
	router.Any("/admin/*w", gin.WrapH(mux))
	//	=== Admin part END ===

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
