package main

import (
	_ "github.com/udistrital/verificacion_firma_digital_mid/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/udistrital/utils_oas/auditoria"
	"github.com/udistrital/utils_oas/xray"
)

func main() {
	AllowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode == "dev" {
		AllowedOrigins = []string{"*"}
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowOrigins: []string{"*"},
		AllowOrigins: AllowedOrigins,
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	xray.InitXRay()
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.Run()
}
