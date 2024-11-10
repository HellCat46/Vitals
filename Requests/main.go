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
	"strconv"
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
	req, err := db.Queryx("SELECT * FROM requests WHERE acceptedBy IS NULL;")
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

func HospitalDeleteRequest(ctx *gin.Context, db *sqlx.DB) {
	reqId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Not a valid Id Query",
		})
		return
	}

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

	_, err = db.NamedExec("DELETE FROM requests WHERE hospitalId = :hospitalId && id = :reqId;", map[string]interface{}{"hospitalId": userId, "reqId": reqId})
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to delete hospital requests",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
	})
}

func DonatorDeleteRequest(ctx *gin.Context, db *sqlx.DB) {
	reqId, err := strconv.Atoi(ctx.Query("id"))
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

	_, err = db.NamedExec("UPDATE requests SET acceptedBy = NULL  WHERE id = :id && acceptedBy = :userId;", map[string]interface{}{"id": reqId, "userId": userId})
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to delete hospital requests",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
	})
}

func AcceptRequest(ctx *gin.Context, db *sqlx.DB) {
	reqId, err := strconv.Atoi(ctx.Query("id"))
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
	tx := db.MustBegin()

	var reqType int
	err = tx.QueryRow(fmt.Sprintf("SELECT type FROM requests WHERE id = %s", reqId)).Scan(&reqType)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to fetch requests data",
		})
		return
	}
	_, err = tx.Exec("UPDATE requests SET acceptedBy = ?  WHERE id = ? && acceptedBy is NULL;", userId, reqId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to update hospital requests",
		})
		tx.Rollback()
		return
	}

	if reqType == 0 {
		_, err = tx.Exec("UPDATE donator SET credits = credits + 10 WHERE userId = ?", userId)
	} else {
		_, err = tx.Exec("UPDATE donator SET credits = credits + 10 WHERE userId = ?", userId)
	}
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to update hospital requests",
		})
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable to update hospital requests",
		})
		tx.Rollback()
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
	})
}
