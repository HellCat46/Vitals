package main

import (
	"Vitals/Auth"
	"Vitals/Db"
	"Vitals/Notification"
	"Vitals/Requests"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
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

	db, err := sqlx.Connect("mysql", "root:Harshit445@tcp(127.0.0.1:3306)/vitals")
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

	r.StaticFile("/", "./vitals-front/dist/index.html")
	r.StaticFS("/assets", http.Dir("./vitals-front/dist/assets"))

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
		Requests.CreateRequest(ctx, db)
	})

	r.GET("/hospital/requests", func(ctx *gin.Context) {
		Requests.HospitalGetRequests(ctx, db)
	})

	r.GET("/donator/requests", func(ctx *gin.Context) {
		Requests.DonatorAllGetRequests(ctx, db)
	})

	r.GET("/donator/accepted_requests", func(ctx *gin.Context) {
		Requests.DonatorAcceptedGetRequests(ctx, db)
	})

	r.GET("/hospital/remRequest", func(ctx *gin.Context) {
		Requests.HospitalDeleteRequest(ctx, db)
	})

	r.GET("/donator/remRequest", func(ctx *gin.Context) {
		Requests.DonatorDeleteRequest(ctx, db)
	})

	r.GET("/donator/acceptRequest", func(ctx *gin.Context) {
		Requests.AcceptRequest(ctx, db)
	})

	r.GET("/donator/getCredits", func(ctx *gin.Context) {
		Requests.GetCredits(ctx, db)
	})

	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
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
