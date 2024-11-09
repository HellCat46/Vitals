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

	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(id) FROM Donator WHERE email = \"%s\"", reqBody.Email)).Scan(&count)
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

	res, err := db.NamedExec("INSERT INTO Donator(name, bloodgroup, email, password, address, phoneno) VALUES(:name, :bloodgroup, :email, :password, :address, :phoneno)",
		map[string]interface{}{"name": reqBody.Name, "bloodgroup": reqBody.BloodGroup, "email": reqBody.Email, "password": string(hashes), "address": reqBody.Address, "phoneno": reqBody.PhoneNo})
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to create an Account.",
		})
		return
	}

	userId, err := res.LastInsertId()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Account Successfully Created but Something went wrong while registering you in.",
		})
		return
	}

	var data []byte
	file, err := reqBody.Image.Open()
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

	err = os.WriteFile(fmt.Sprintf("Data/Donar/%d.png", userId), data, 667)
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

type HospitalRegBody struct {
	Name        string                `form:"name"`
	Email       string                `form:"email"`
	Password    string                `form:"password"`
	Address     string                `form:"address"`
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
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(id) FROM Donator WHERE email = \"%s\"", reqBody.Email)).Scan(&count)
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

	res, err := db.NamedExec("INSERT INTO hospital(name, email, password, address, phoneno) VALUES(:name, :email, :password, :address, :phoneno)",
		map[string]interface{}{"name": reqBody.Name, "email": reqBody.Email, "password": string(hashes), "address": reqBody.Address, "phoneno": reqBody.PhoneNo})
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unable to create an Account.",
		})
		return
	}

	userId, err := res.LastInsertId()
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Account Successfully Created but Something went wrong while registering you in.",
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
