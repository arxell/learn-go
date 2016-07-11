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
)

var DB = make(map[string]string)

type OrderItem struct {
	ID       uint
	Order_id int
	Status   string
}

func (OrderItem) TableName() string {
	return "stat_orderitem"
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
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		ois := []OrderItem{}
		DB.Limit(10).Find(&ois)
		fmt.Println(ois)
		c.String(200, "pong")
	})

	// Initalize
	Admin := admin.New(&qor.Config{DB: DB})

	// Create resources from GORM-backend model
	Admin.AddResource(&OrderItem{})

	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)

	r.Any("/admin/*w", gin.WrapH(mux))
	r.Run(":8080")
}
