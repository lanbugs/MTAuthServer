package mtauthserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
)

// @Summary Authentication
// @Description Authenticate to get JWT token
// @Produce json
// @Param UsernamePassword body UsernamePassword true "Username and Password"
// @Success 200 {object} ResponseToken "token response"
// @failure 400 {object} ResponseAuthError "error response"
// @failure 401 {object} ResponseAuthError "error response"
// @Router /auth [post]
func Auth(c *gin.Context) {
	var data UsernamePassword

	if err := c.BindJSON(&data); err != nil {
		log.Errorf("error dataformat not correct")
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "data format not correct"})
		return
	}

	// Check basics
	if len(data.Username) < 5 || len(data.Username) > 50 {
		log.Errorf("error username too short or too long min. 5 max. 50")
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "error username too short or too long min. 5 max. 50"})
		return
	}

	if len(data.Password) < 5 || len(data.Password) > 128 {
		log.Errorf("error password too short or too long min. 5 max. 128")
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "error password too short or too long min. 5 max. 128"})
		return
	}

	l := ConnectLDAP()

	if checkAuthentication(l, data.Username, data.Password) {
		groups := getGroupsofUser(l, data.Username)
		token := generate_token(data.Username, groups)

		c.JSON(http.StatusOK, gin.H{"status": "ok", "username": data.Username, "groups": groups, "token": token})
		return
	} else {

		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": "Authorization failed"})
		return
	}

}

// @Summary Introspect
// @Description Check JWT token
// @Produce json
// @Param Token body Token true "Token"
// @failure 400 {object} ResponseAuthError "error response"
// @failure 401 {object} ResponseAuthError "error response"
// @Success 200 {object} ResponseVerify "verify response"
// @Router /introspect [post]
func Introspect(c *gin.Context) {
	var data Token

	if err := c.BindJSON(&data); err != nil {
		log.Errorf("data format not correct")
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "data format not correct"})
		return
	}

	var valid bool = false
	var username string = ""
	var groups string = ""

	valid, username, groups = validate_token(data.Token)

	if valid {
		g := []string{}

		if err := json.Unmarshal([]byte(groups), &g); err != nil {
			fmt.Printf("Error: %v", err)
		}
		c.JSON(http.StatusOK, gin.H{"status": "valid", "username": username, "groups": g})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": "token invalid"})
		return
	}

}

// @Summary Verification
// @Description Verify JWT token
// @Produce json
// @Param app_name path string true "application name"
// @Param Authorization header string true "authentication token"
// @Success 200 {object} ResponseVerify "verify response"
// @failure 400 {object} ResponseError "error response"
// @failure 401 {object} ResponseError "error response"
// @Router /verify/{app_name} [get]
func VerifyToken(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	app_name := c.Param("app_name")

	app_name_valid := regexp.MustCompile(`^[A-Za-z_-]+$`).MatchString(app_name)

	if !app_name_valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "app_name has invalid chars in it, only A-Z a-z and _-."})
		return
	}

	// Check authorization header
	if authHeader == "" {
		c.JSON(401, gin.H{"message": "Authorization header missing."})
		return
	}

	// Check format of Token
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		c.JSON(401, gin.H{"message": "Authorization header format not valid."})
		return
	}

	// Extrahiere den Token
	token := authHeaderParts[1]

	var verify bool = false
	var username string = ""
	var groups string = ""

	verify, username, groups = validate_token(token)

	if verify {

		g := []string{}

		if err := json.Unmarshal([]byte(groups), &g); err != nil {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "unable to load groups from jwt token."})
		}

		log.WithFields(log.Fields{"groups": g, "username": username, "app_name": app_name}).Info("Authorization successful.")
		c.JSON(http.StatusOK, gin.H{"status": "valid", "username": username, "groups": g, "app_name": app_name})
		return
	} else {
		log.WithFields(log.Fields{"app_name": app_name}).Warnf("Authorization failed.")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": "token invalid", "app_name": app_name})
		return
	}

}
