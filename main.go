package main

import (
	docs "MTAuthServer/docs"
	"MTAuthServer/mtauthserver"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/toorop/gin-logrus"
	"net/http"
	"os"
	"syscall"
)

const (
	NAME        = "MTAuthServer"
	DESCRIPTION = "JWT authentication server to Active Directory"
	VERSION     = "1.0.3"
)

func Help() {
	fmt.Printf("\n%v - %v - Version %v\n", NAME, DESCRIPTION, VERSION)
	fmt.Println("Written by Maximilian Thoma 2024")

	flag.Usage()
	syscall.Exit(0)
}

func main() {
	log := logrus.New()

	var help *bool = flag.BoolP("help", "h", false, "Show help")

	if *help {
		Help()
	}

	cnf := mtauthserver.LoadConfig()

	if cnf.LogtoFile {
		f, err := os.OpenFile(cnf.LogFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Printf("unable to open logfile error: %v", err)
		}

		log.SetOutput(f)
	}

	if cnf.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(ginlogrus.Logger(log), gin.Recovery())
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = NAME
	docs.SwaggerInfo.Version = VERSION
	docs.SwaggerInfo.Description = "JWT Authentication to AD"
	api := r.Group("/api/v1")
	api.POST("/auth", mtauthserver.Auth)
	api.POST("/introspect", mtauthserver.Introspect)
	api.GET("/verify/:app_name", mtauthserver.VerifyToken)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.Run()

}
