package Auth

import (
	"Vitals/Db"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"io"
	"strings"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DonatorLogin(ctx *gin.Context, db *sqlx.DB) {
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Unparsable Body",
		})
		return
	}

	var loginData LoginData
	err = json.Unmarshal(data, &loginData)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Unparsable Body",
		})
		return
	}

	if len(loginData.Email) == 0 || len(loginData.Password) == 0 {
		ctx.JSON(422, map[string]string{
			"error": "One of the required parameters are missing",
		})
		return
	}

	var user Db.User
	rows, err := db.Queryx(fmt.Sprintf("SELECT * FROM user WHERE email = '%s';", loginData.Email))
	if rows == nil {
		ctx.JSON(404, map[string]string{
			"error": "User does not exist",
		})
		return
	}
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unexpected error",
		})
		return
	}
	if rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			println(err.Error())
			ctx.JSON(500, map[string]string{
				"error": "Unexpected error",
			})
			return
		}
	}

	println(user.Id)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		ctx.JSON(403, map[string]string{
			"error": "Invalid Credentials",
		})
		return
	}
	//
	//if err != nil {
	//	println(err.Error())
	//	ctx.JSON(422, map[string]string{
	//		"error": "Unable to Confirm Account Info. Please try again later.",
	//	})
	//	return
	//}

	token, err := GenerateToken(int64(user.Id))
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Account Successfully Created but Something went wrong while logging you in.",
		})
		return
	}

	ctx.JSON(200, map[string]string{
		"token": token,
	})
}

func HospitalLogin(ctx *gin.Context, db *sqlx.DB) {
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Unparsable Body",
		})
		return
	}

	var loginData LoginData
	err = json.Unmarshal(data, &loginData)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Unparsable Body",
		})
		return
	}

	if len(loginData.Email) == 0 || len(loginData.Password) == 0 {
		ctx.JSON(422, map[string]string{
			"error": "One of the required parameters are missing",
		})
		return
	}

	var user Db.User
	rows, err := db.Queryx(fmt.Sprintf("SELECT * FROM user WHERE email = '%s';", loginData.Email))
	if rows == nil {
		ctx.JSON(404, map[string]string{
			"error": "User does not exist",
		})
		return
	}
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Unexpected error",
		})
		return
	}
	if rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			println(err.Error())
			ctx.JSON(500, map[string]string{
				"error": "Unexpected error",
			})
			return
		}
	}

	println(user.Id)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		ctx.JSON(403, map[string]string{
			"error": "Invalid Credentials",
		})
		return
	}
	//
	//if err != nil {
	//	println(err.Error())
	//	ctx.JSON(422, map[string]string{
	//		"error": "Unable to Confirm Account Info. Please try again later.",
	//	})
	//	return
	//}

	token, err := GenerateToken(int64(user.Id))
	if err != nil {
		println(err.Error())
		ctx.JSON(500, map[string]string{
			"error": "Account Successfully Created but Something went wrong while logging you in.",
		})
		return
	}

	ctx.JSON(200, map[string]string{
		"token": token,
	})
}

func DecodeUnsignedJWT(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 && len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("error decoding payload: %v", err)
	}

	// Parse the JSON
	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return "", fmt.Errorf("error parsing claims: %v", err)
	}

	// Extract the "sub" claim
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub claim not found or not a string")
	}

	return sub, nil
}
