package mtauthserver

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

func generate_token(username string, groups []string) string {
	cnf := LoadConfig()

	json_groups, err := json.Marshal(groups)

	if err != nil {
		log.Errorf("Error convert to JSON: %v\n", err)
	}

	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["groups"] = string(json_groups)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(cnf.SecretKey)

	signedToken, err := token.SignedString(secretKey)

	if err != nil {
		log.Errorf("Error creating token: %v", err)
	}

	return signedToken
}

func validate_token(token string) (valid bool, username string, groups string) {
	cnf := LoadConfig()

	secretKey := []byte(cnf.SecretKey)

	tokenx, err := jwt.Parse(token, func(tokenx *jwt.Token) (interface{}, error) {
		// Check signature algorithm
		if _, ok := tokenx.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signature algorithm: %v", tokenx.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Errorf("Error token validation: %v\n", err)
		return false, "", ""
	}

	// Check token if valid
	if claims, ok := tokenx.Claims.(jwt.MapClaims); ok && tokenx.Valid {
		// Check expiry
		expTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expTime) {
			log.Errorf("Token expired.")
			return false, "", ""
		}

		// Das Token ist gültig und nicht abgelaufen
		log.WithFields(log.Fields{"username": claims["username"], "groups": claims["groups"]}).Infoln("Token valid.")

		return true, claims["username"].(string), claims["groups"].(string)
	} else {
		log.Infoln("Token is invalid.")
		return false, "", ""
	}

}
