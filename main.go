package main

import (
	"Vitals/Auth"
	"Vitals/Db"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
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

	r := gin.Default()
	r.POST("/donator/register", func(ctx *gin.Context) {
		Auth.DonatorReg(ctx, db)
	})
	r.POST("/hospital/register", func(ctx *gin.Context) {
		Auth.HospitalReg(ctx, db)
	})

	r.POST("/donator/login", func(ctx *gin.Context) {
		Auth.DonatorReg(ctx, db)
	})
	r.POST("/hospital/login", func(ctx *gin.Context) {
		Auth.HospitalLogin(ctx, db)
	})

	r.Run()
}

// re_T1c6Yzeq_4qyVLVoMFxbMGNKB2rKrZqJF
