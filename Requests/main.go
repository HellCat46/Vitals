package Requests

import (
	"Vitals/Auth"
	"Vitals/Db"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CreateReqBody struct {
	BloodGroup string `json:"blood_group"`
	NeedType   string `json:"need_type"`
}

func CreateRequest(ctx *gin.Context, db *sqlx.DB) {
	token := ctx.Request.Header.Get("X-TOKEN")
	if token == "" {
		ctx.JSON(401, gin.H{
			"error": "You need to be logged in to perform this action",
		})
		return
	}

	userId, err := Auth.DecodeUnsignedJWT(token)
	if err != nil {
		ctx.JSON(401, gin.H{
			"error": "Invalid token",
		})
		return
	}

	res, err := db.Queryx("SELECT * FROM hospital WHERE userId = :userId", userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to fetch hospital data",
		})
		return
	}

	var hospitalData Db.Hospital
	if res.Next() {
		err := res.StructScan(&hospitalData)
		if err != nil {
			ctx.JSON(500, gin.H{
				"" +
					"error": "Unable to fetch hospital data",
			})
			return
		}
	} else {
		ctx.JSON(500, gin.H{
			"error": "Unable to fetch hospital data",
		})
		return
	}

}

type ReqBody struct {
	Numbers    []string `json:"numbers"`
	Hospital   string   `json:"hospital"`
	Addr       string   `json:"addr"`
	BloodGroup string   `json:"blood_group"`
}
