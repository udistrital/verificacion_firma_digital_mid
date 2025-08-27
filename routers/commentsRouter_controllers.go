package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/verificacion_firma_digital_mid/controllers:VerificarFirmaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/verificacion_firma_digital_mid/controllers:VerificarFirmaController"],
        beego.ControllerComments{
            Method: "PostVerificarFirma",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
