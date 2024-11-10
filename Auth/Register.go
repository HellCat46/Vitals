package Auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"os"
	"time"
)

type DonatorRegBody struct {
	Name       string                `form:"name"`
	BloodGroup string                `form:"blood_group"`
	Email      string                `form:"email"`
	Password   string                `form:"password"`
	Address    string                `form:"address"`
	PhoneNo    string                `form:"phoneno"`
	Pincode    string                `form:"pincode"`
	Image      *multipart.FileHeader `form:"image"`
}

func DonatorReg(ctx *gin.Context, db *sqlx.DB) {
	var reqBody DonatorRegBody
	err := ctx.Bind(&reqBody)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Body"})
		return
	}

	var count int

	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(id) FROM user WHERE email = \"%s\"", reqBody.Email)).Scan(&count)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Failed to Process the Request",
		})
		return
	}

	if count > 0 {
		ctx.JSON(409, map[string]string{
			"error": "User Already Exists",
		})
		return
	}

	hashes, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.MinCost)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the password",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to start a transaction",
		})
		return
	}

	res, err := tx.Exec("INSERT INTO user(email, password, type) VALUES(?, ?, ?);", reqBody.Email, string(hashes), 1)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the request",
		})
		tx.Rollback()
		return
	}

	userId, err := res.LastInsertId()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the request",
		})
		tx.Rollback()
		return
	}

	res, err = tx.Exec("INSERT INTO donator(userId, name, bloodgroup, address, phoneno, pincode) VALUES(?, ?, ?, ?, ?, ?)", userId, reqBody.Name, reqBody.BloodGroup, reqBody.Address, reqBody.PhoneNo, reqBody.Pincode)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to create an Account.",
		})
		tx.Rollback()
		return
	}

	var data []byte
	file, err := reqBody.Image.Open()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the ID Proof",
		})
		tx.Rollback()
		return

	}
	_, err = file.Read(data)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the ID Proof",
		})
		tx.Rollback()
		return
	}

	err = os.WriteFile(fmt.Sprintf("Data/Donor/%d.png", userId), data, 667)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to save the ID Proof",
		})
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to commit the transaction",
		})
		tx.Rollback()
		return
	}

	token, err := GenerateToken(userId)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Account Successfully Created but Something went wrong while registering you in.",
		})
		return
	}

	ctx.JSON(200, map[string]string{
		"token": token,
	})
}

type HospitalRegBody struct {
	Name        string                `form:"name"`
	Email       string                `form:"email"`
	Password    string                `form:"password"`
	Address     string                `form:"address"`
	Pincode     string                `form:"pincode"`
	PhoneNo     string                `form:"phoneno"`
	Certificate *multipart.FileHeader `form:"cert"`
}

func HospitalReg(ctx *gin.Context, db *sqlx.DB) {
	var reqBody HospitalRegBody
	err := ctx.Bind(&reqBody)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Body"})
		return
	}

	var count int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(id) FROM user WHERE email = \"%s\"", reqBody.Email)).Scan(&count)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Failed to Process the Request",
		})
		return
	}

	if count > 0 {
		ctx.JSON(409, map[string]string{
			"error": "User Already Exists",
		})
		return
	}

	hashes, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.MinCost)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the password",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to start a transaction",
		})
		return
	}

	res, err := tx.Exec("INSERT INTO user(email, password, type) VALUES(?, ?, ?);", reqBody.Email, string(hashes), 0)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the request",
		})
		tx.Rollback()
		return
	}

	userId, err := res.LastInsertId()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the request",
		})
		tx.Rollback()
		return
	}

	res, err = tx.Exec("INSERT INTO hospital(userId,name, address, phoneno, pincode) VALUES(?, ?, ?, ?, ?)", userId, reqBody.Name, reqBody.Address, reqBody.PhoneNo, reqBody.Pincode)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to create an Account.",
		})
		return
	}

	var data []byte
	file, err := reqBody.Certificate.Open()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the ID Proof",
		})
		return

	}
	_, err = file.Read(data)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to process the ID Proof",
		})
		return
	}

	err = os.WriteFile(fmt.Sprintf("Data/Hospital/%d.png", userId), data, 667)
	if err != nil {
		ctx.JSON(500, map[string]string{
			"error": "Unable to save the ID Proof",
		})
		return
	}

	token, err := GenerateToken(userId)
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Account Successfully Created but Something went wrong while registering you in.",
		})
		return
	}

	ctx.JSON(200, map[string]string{
		"token": token,
	})
}

func GenerateToken(userId int64) (token string, err error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userId),
		"iss": "vitals",
		"exp": time.Now().Add(time.Hour * 720).Unix(),
		"iat": time.Now().Unix(),
	})

	return claims.SigningString()
}
