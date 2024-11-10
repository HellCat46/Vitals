package Requests

import (
	"Vitals/Auth"
	"Vitals/Db"
	"Vitals/Notification"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"io"
)

type CreateReqBody struct {
	BloodGroup string `json:"blood_group"`
	NeedType   int    `json:"need_type"`
	Unit       int    `json:"unit"`
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

	println(userId)
	res, err := db.Queryx(fmt.Sprintf("SELECT * FROM hospital WHERE userId = %s", userId))
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
			println(err.Error())
			ctx.JSON(500, gin.H{
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

	_, err = db.NamedExec("INSERT INTO requests(hospitalId, type, bloodgroup, unit) VALUES(:hosId, :type, :bloodgroup, :unit)",
		map[string]interface{}{
			"hosId":      hospitalData.UserId,
			"type":       createReqBody.NeedType,
			"bloodgroup": createReqBody.BloodGroup,
			"unit":       createReqBody.Unit,
		})
	if err != nil {
		println(err.Error())
		ctx.JSON(500, gin.H{
			"error": "Unable to create a request",
		})
		return
	}

	if createReqBody.NeedType == 0 {

		res, err := db.Queryx("SELECT name, phoneno FROM donator WHERE pincode = :pincode && bloodgroup = : bloodgroup", hospitalData.Pincode, createReqBody.BloodGroup)
		if err == nil {
			type Donator struct {
				Name    string `db:"name"`
				Phoneno string `db:"phoneno"`
			}
			var users []Notification.User
			for res.Next() {
				var donator Donator
				err := res.StructScan(&donator)
				if err != nil {
					continue
				}
				users = append(users, Notification.User{
					Number: donator.Phoneno,
					Name:   donator.Name,
				})
			}

			if len(users) != 0 {
				err := Notification.SendBulkMessage(Notification.ReqBody{
					Users:      users,
					Hospital:   hospitalData.Name,
					Addr:       hospitalData.Address,
					BloodGroup: createReqBody.BloodGroup,
					Type:       0,
				})
				if err != nil {
					println(err.Error())
				} else {
					ctx.JSON(200, map[string]string{
						"status": "Success",
					})
					return
				}
			}
		}
	}

	res, err = db.Queryx(fmt.Sprintf("SELECT name, phoneno FROM donator WHERE bloodgroup = '%s'", createReqBody.BloodGroup))
	if err != nil {
		println(err.Error())
		ctx.JSON(200, gin.H{
			"error": "Request was successfully Created but Unable to inform users.",
		})
		return
	}
	type Donator struct {
		Name    string `db:"name"`
		Phoneno string `db:"phoneno"`
	}
	var users []Notification.User
	for res.Next() {
		var donator Donator
		err := res.StructScan(&donator)
		if err != nil {
			continue
		}
		users = append(users, Notification.User{
			Number: donator.Phoneno,
			Name:   donator.Name,
		})
	}

	println(len(users))
	if len(users) != 0 {
		err := Notification.SendBulkMessage(Notification.ReqBody{
			Users:      users,
			Hospital:   hospitalData.Name,
			Addr:       hospitalData.Address,
			BloodGroup: createReqBody.BloodGroup,
			Type:       1,
		})
		if err != nil {
			println(err.Error())
		}
	}

	ctx.JSON(200, map[string]string{
		"status": "Success",
	})
}

func HospitalGetRequests(ctx *gin.Context, db *sqlx.DB) {
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

	println(userId)
	req, err := db.Queryx(fmt.Sprintf("SELECT * FROM requests WHERE hospitalId = %s;", userId))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to fetch blood requests",
		})
		return
	}

	var requests []Db.Request
	for req.Next() {
		var request Db.Request
		err := req.StructScan(&request)
		if err != nil {
			println(err.Error())
			continue
		}

		requests = append(requests, request)
	}

	ctx.JSON(200, gin.H{
		"requests": requests,
	})
}

func DonatorAcceptedGetRequests(ctx *gin.Context, db *sqlx.DB) {
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

	println(userId)
	req, err := db.Queryx(fmt.Sprintf("SELECT * FROM requests WHERE acceptedBy = %s;", userId))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to fetch blood requests",
		})
		return
	}

	var requests []Db.Request
	for req.Next() {
		var request Db.Request
		err := req.StructScan(&request)
		if err != nil {
			println(err.Error())
			continue
		}

		requests = append(requests, request)
	}

	ctx.JSON(200, gin.H{
		"requests": requests,
	})
}

func DonatorAllGetRequests(ctx *gin.Context, db *sqlx.DB) {
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

	println(userId)
	req, err := db.Queryx(fmt.Sprintf("SELECT * FROM requests WHERE acceptedBy != %s;", userId))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to fetch blood requests",
		})
		return
	}

	var requests []Db.Request
	for req.Next() {
		var request Db.Request
		err := req.StructScan(&request)
		if err != nil {
			println(err.Error())
			continue
		}

		requests = append(requests, request)
	}

	ctx.JSON(200, gin.H{
		"requests": requests,
	})
}

func DeleteRequest(ctx *gin.Context, db *sqlx.DB) {
	token := ctx.Request.Header.Get("X-TOKEN")
	if token == "" {
		ctx.JSON(401, gin.H{
			"error": "You need to be logged in to perform this action",
		})
		return
	}
}

func AcceptRequest(ctx *gin.Context, db *sqlx.DB) {
	token := ctx.Request.Header.Get("X-TOKEN")
	if token == "" {
		ctx.JSON(401, gin.H{
			"error": "You need to be logged in to perform this action",
		})
		return
	}
}
