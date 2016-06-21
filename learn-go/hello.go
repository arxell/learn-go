package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"log"
	"net/http"
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

func main() {
	DB, err := gorm.Open("postgres", "host= port=5432 user=postgres dbname= password=")

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
