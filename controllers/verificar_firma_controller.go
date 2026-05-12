package controllers

import (
	"github.com/astaxie/beego"
	"log"
	"encoding/json"
	"github.com/udistrital/verificacion_firma_digital_mid/services"
	"github.com/udistrital/verificacion_firma_digital_mid/models"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
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
	log.Println("[trace] verificar_firma.http.enter")
	defer errorhandler.HandlePanic(&c.Controller)

	var archivos []models.EmailAttachment
	log.Println("[trace] verificar_firma.http.unmarshal.start")
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &archivos); err != nil {
		log.Println("[trace] verificar_firma.http.unmarshal.fail")
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar el cuerpo de la solicitud: "+err.Error())
		c.ServeJSON()
		return
	}

	log.Println("[trace] verificar_firma.http.unmarshal.ok")
	if len(archivos) == 0 {
		log.Println("[trace] verificar_firma.http.empty_array")
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "El array de archivos está vacío")
		c.ServeJSON()
		return
	}

	archivo := archivos[0]

	// Validaciones de campos
	log.Println("[trace] verificar_firma.http.input.loaded")
	if archivo.PdfBase64 == "" || archivo.Firma == "" || archivo.UrlFileUp == "" {
		log.Println("[trace] verificar_firma.http.validation.fail")
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Los campos pdf_base64, firma y urlFileUp son obligatorios")
		c.ServeJSON()
		return
	}

	// Verificar PDF con ClamAV
	log.Println("[trace] verificar_firma.clamav.start")
	resultadoClamAV := services.VerificarPDFBase64(archivo.PdfBase64)
	log.Printf("[trace] verificar_firma.clamav.end | success=%v status=%d\n", resultadoClamAV.Success, resultadoClamAV.Status)
	if !resultadoClamAV.Success {
		c.Ctx.Output.SetStatus(resultadoClamAV.Status)
		c.Data["json"] = resultadoClamAV
		c.ServeJSON()
		return
	}
	
	// Obtener resultado del antivirus y formatearlo
	var virusResult map[string]interface{}
	if clamAVData, ok := resultadoClamAV.Data.(map[string]interface{}); ok {
		if virusData, ok := clamAVData["Virus"].(map[string]interface{}); ok {
			virusResult = virusData
		} else {
			// fallback si viene sin clave "Virus"
			virusResult = map[string]interface{}{
				"message":    "Verificación de virus completada correctamente.",
				"archive":    clamAVData["status"],
				"statusCode": 200,
			}
		}
	} else {
		virusResult = map[string]interface{}{
			"message":    "Error al procesar el archivo",
			"archive":    "unknown",
			"statusCode": 500,
		}
	}
	
	log.Println("[trace] verificar_firma.signature.start")
	token := c.Ctx.Input.Header("Authorization")
	// Verificar firma electrónica
	respuestaFirma := services.PostVerificarFirma(archivo, token)
	log.Printf("[trace] verificar_firma.signature.end | success=%v status=%d\n", respuestaFirma.Success, respuestaFirma.Status)

	// Unificar resultados
	if !respuestaFirma.Success {
		c.Ctx.Output.SetStatus(respuestaFirma.Status)
		c.Data["json"] = respuestaFirma
		c.ServeJSON()
		return
	}

	// Obtener resultado de firma
	var firmaResult map[string]interface{}
	var fileEqual bool

	if dataMap, ok := respuestaFirma.Data.(map[string]interface{}); ok {
		if verificacion, exists := dataMap["Verificacion"].(map[string]interface{}); exists {
			firmaResult = verificacion
			if fe, ok := verificacion["fileEqual"].(bool); ok {
				fileEqual = fe
			}
		}
	}

	// Mensaje final combinado
	mensajeFinal := "El archivo PDF está limpio."
	if virusResult["archive"] == "infected" {
		mensajeFinal = "El archivo contiene virus."
	} else if !fileEqual {
		mensajeFinal += " Pero los archivos no coinciden con la firma proporcionada."
	}

	// Estructura final de respuesta
	dataFinal := map[string]interface{}{
		"Virus":       virusResult,
		"Verificacion": firmaResult,
	}

	respuestaFinal := requestresponse.APIResponseDTO(
		true,
		200,
		dataFinal,
		mensajeFinal,
	)

	log.Println("[trace] verificar_firma.http.response.200")
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = respuestaFinal
	c.ServeJSON()

}
