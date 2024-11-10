package Requests

import (
	"Vitals/Auth"
	"Vitals/Db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"io"
)

type CreateReqBody struct {
	BloodGroup string `json:"blood_group"`
	NeedType   int    `json:"need_type"`
}

func CreateRequest(ctx *gin.Context, db *sqlx.DB) {
	token := ctx.Request.Header.Get("X-TOKEN")
	if token == "" {
		ctx.JSON(401, gin.H{
			"error": "You need to be logged in to perform this action",
		})
		return
	}

	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Unparsable request body",
		})
		return
	}
	var createReqBody CreateReqBody
	err = json.Unmarshal(data, &createReqBody)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Unparsable request body",
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
	if createReqBody.NeedType == 0 {
		res, err := db.Queryx("SELECT name, phoneno FROM donator WHERE pincode = :pincode && bloodgroup = : bloodgroup", hospitalData.Pincode, createReqBody.BloodGroup)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": "Unable to fetch donor's data",
			})
		}

		type Donator struct {
			Name    string `db:"name"`
			Phoneno string `db:"phoneno"`
		}
		var donators []Donator
		for res.Next() {
			var donator Donator
			err := res.StructScan(&donator)
			if err != nil {
				continue
			}
			donators = append(donators, donator)
		}
	}

}

func GetRequests(ctx *gin.Context, db *sqlx.DB) {
	token := ctx.Request.Header.Get("X-TOKEN")
	if token == "" {
		ctx.JSON(401, gin.H{})
	}
}

func AcceptRequest(ctx *gin.Context, db *sqlx.DB) {
	token := ctx.Request.Header.Get("X-TOKEN")
	if token == "" {
		ctx.JSON(401, gin.H{})
	}
}
