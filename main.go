package main

import (
	"Vitals/Auth"
	"Vitals/Db"
	"Vitals/Notification"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"os"
	"time"
)

func main() {
	err := os.MkdirAll("Data/Donor", 0750)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll("Data/Hospital", 0750)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("mysql", "root:kunal@tcp(127.0.0.1:3306)/vitals")
	if err != nil {
		panic(err)
	}

	for idx, schema := range Db.Schemas {
		db.MustExec(schema)
		fmt.Printf("Successfully created table no %d.\n", idx+1)
	}

	cronManager := cron.New()

	cronManager.AddFunc("*/1 * * * *", func() {
		err := Notification.SendWarning("8302071621")
		if err != nil {
			print(err.Error())
		}
	})

	r := gin.Default()
	r.POST("/donator/register", func(ctx *gin.Context) {
		Auth.DonatorReg(ctx, db)
	})
	r.POST("/hospital/register", func(ctx *gin.Context) {
		Auth.HospitalReg(ctx, db)
	})

	r.POST("/donator/login", func(ctx *gin.Context) {
		Auth.DonatorLogin(ctx, db)
	})
	r.POST("/hospital/login", func(ctx *gin.Context) {
		Auth.HospitalLogin(ctx, db)
	})

	r.POST("/hospital/createRequest", func(ctx *gin.Context) {

	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Run()
	cronManager.Run()
}

// re_T1c6Yzeq_4qyVLVoMFxbMGNKB2rKrZqJF
