package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/udistrital/verificacion_firma_digital_mid/services"
	"github.com/udistrital/verificacion_firma_digital_mid/models"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
	"fmt"
)

// VerificarFirmaController operations for VerficarFirma
type VerificarFirmaController struct {
	beego.Controller
}

// URLMapping ...
func (c *VerificarFirmaController) URLMapping() {
	c.Mapping("PostVerificarFirma", c.PostVerificarFirma)
}

// PostVerificarFirma ...
// @Title PostVerificarFirma
// @Description Verifica si un PDF base64 está limpio usando ClamAV
// @Param	body		body 	models.EmailAttachment	true		"Base64 del PDF y hash"
// @Success 200 {object} map[string]interface{}
// @Failure 404 body is empty
// @router "/" [post]
func (c *VerificarFirmaController) PostVerificarFirma() {
	defer errorhandler.HandlePanic(&c.Controller)

	var archivos []models.EmailAttachment
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &archivos); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar el cuerpo de la solicitud: "+err.Error())
		c.ServeJSON()
		return
	}

	if len(archivos) == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "El array de archivos está vacío")
		c.ServeJSON()
		return
	}

	archivo := archivos[0]

	// Validaciones de campos
	if archivo.PdfBase64 == "" || archivo.Firma == "" || archivo.UrlFileUp == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Los campos pdf_base64, firma y urlFileUp son obligatorios")
		c.ServeJSON()
		return
	}

	// Verificar PDF con ClamAV
	resultadoClamAV := services.VerificarPDFBase64(archivo.PdfBase64)
	if !resultadoClamAV.Success {
		c.Ctx.Output.SetStatus(resultadoClamAV.Status)
		c.Data["json"] = resultadoClamAV
		c.ServeJSON()
		return
	}
	fmt.Println("Resultado de ClamAV:", resultadoClamAV)

	token := c.Ctx.Input.Header("Authorization")
	// Verificar firma electrónica
	respuestaFirma := services.PostVerificarFirma(archivo, token)

	// Unificar resultados
	if !respuestaFirma.Success {
		c.Ctx.Output.SetStatus(respuestaFirma.Status)
		c.Data["json"] = respuestaFirma
		c.ServeJSON()
		return
	}

	// Extraer valor de fileEqual del resultado
	dataMap, ok := respuestaFirma.Data.(map[string]interface{})
	fmt.Println("DataMap:", dataMap)
	var fileEqual bool
	if ok {
		if verificacion, exists := dataMap["Verificacion"].(map[string]interface{}); exists {
			if fe, ok := verificacion["fileEqual"].(bool); ok {
				fileEqual = fe
			}
		}
	}
	
	mensajeFinal := "El archivo PDF está limpio."
	if !fileEqual {
		mensajeFinal += " Pero los archivos no coinciden con la firma proporcionada."
	}
	fmt.Println("fileEqual:", fileEqual)

	// Construir respuesta final
	respuestaFinal := requestresponse.APIResponseDTO(
		true,
		200,
		respuestaFirma.Data,
		mensajeFinal,
	)

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = respuestaFinal
	c.ServeJSON()

}


/*func (c *VerificarFirmaController) PostVerificarFirma() {
	defer errorhandler.HandlePanic(&c.Controller)

	var archivo models.EmailAttachment
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &archivo); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar el cuerpo de la solicitud: "+err.Error())
		c.ServeJSON()
		return
	}

	if archivo.PdfBase64 == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "El campo pdf_base64 no puede estar vacío")
		c.ServeJSON()
		return
	}

	// ✅ Solo un valor de retorno
	respuesta := services.VerificarPDFBase64(archivo.PdfBase64)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}*/
