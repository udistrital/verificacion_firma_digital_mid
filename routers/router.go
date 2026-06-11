// @APIVersion 1.0.0
// @Title Verficar Firma MID - Verficar Firma Digital
// @Description Microservicio MID de Verficar Firma MID que complementa Verficar Firma Digital
package routers

import (
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/verificacion_firma_digital_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.ErrorController(&errorhandler.ErrorHandlerController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/verificar_firma",
			beego.NSInclude(
				&controllers.VerificarFirmaController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
