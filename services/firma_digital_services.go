/*package services

import (
	"os"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/verificacion_firma_digital_mid/models"
)

func PostVerificarFirma(archivo models.EmailAttachment, token string) requestresponse.APIResponse {
	log.Println("[trace] service.signature.start")
	log.Println("[trace] service.signature.start")
	// Establecer token como variable de entorno para el paquete request
	os.Setenv("Authorization", token)
	fmt.Println("Token establecido:", token)

	url := beego.AppConfig.String("FirmaElectronicaService") + "verify"

	payload := []models.EmailAttachment{archivo}

	log.Println("[trace] service.signature.response.read.ok")
	var resultado map[string]interface{}

	err := request.SendJson(url, "POST", &resultado, payload)
	if err != nil {
		beego.Error("Error en POST a verificador firma electrónica:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error en petición POST: %v", err))
	}

	return requestresponse.APIResponseDTO(true, 200, resultado, "Verificación de firma completada correctamente.")
}*/


package services

import (
	"bytes"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/astaxie/beego"
	"github.com/udistrital/verificacion_firma_digital_mid/models"
	"github.com/udistrital/utils_oas/requestresponse"
)

func PostVerificarFirma(archivo models.EmailAttachment, token string) requestresponse.APIResponse {
	log.Println("[trace] service.signature.start")
	url := beego.AppConfig.String("FirmaElectronicaService") + "verify"

	payload := []models.PayloadVerificacion{
		{
			FileUp:    archivo.PdfBase64,
			Firma:     archivo.Firma,
			UrlFileUp: archivo.UrlFileUp,
		},
	}

	// Serializar el payload a JSON
	log.Println("[trace] service.signature.marshal.start")
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		beego.Error("Error al serializar JSON:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al serializar JSON: "+err.Error())
	}
	log.Println("[trace] service.signature.marshal.ok")

	// Crear la solicitud HTTP con body y header Authorization
	log.Println("[trace] service.signature.request.build.start")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		beego.Error("Error al crear solicitud HTTP:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al crear solicitud HTTP: "+err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	//req.Header.Add("Authorization", "Bearer 299a6897-f955-39c1-b4ce-9fd731387b3d") // Usar Add para evitar reemplazar si ya existe
	log.Println("[trace] service.signature.request.build.ok")
	//fmt.Println("Token quemado:", "Bearer 299a6897-f955-39c1-b4ce-9fd731387b3d")

	client := &http.Client{}
	log.Println("[trace] service.signature.request.send.start")
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Error al enviar solicitud HTTP:", err)

		verificacion := map[string]interface{}{
			"statusCode": 500,
			"fileEqual":  false,
			"Message":    "Error al enviar solicitud HTTP: " + err.Error(),
		}

		return requestresponse.APIResponseDTO(false, 500, map[string]interface{}{"Verificacion": verificacion}, verificacion["Message"].(string))
	}
	log.Printf("[trace] service.signature.request.send.end | status=%d\n", resp.StatusCode)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var resultado map[string]interface{}
	log.Println("[trace] service.signature.response.unmarshal.start")
	if err := json.Unmarshal(body, &resultado); err != nil {
		log.Println("[trace] service.signature.response.unmarshal.fail")
		beego.Error("Error al parsear respuesta JSON:", err)

		verificacion := map[string]interface{}{
			"statusCode": resp.StatusCode,
			"fileEqual":  false,
			"Message":    "Error al parsear respuesta JSON: " + err.Error(),
		}

		return requestresponse.APIResponseDTO(false, 500, map[string]interface{}{"Verificacion": verificacion}, verificacion["Message"].(string))
	}

	var fileEqual bool
	if data, ok := resultado["res"].([]interface{}); ok && len(data) > 0 {
		if item, ok := data[0].(map[string]interface{}); ok {
			fileEqual, _ = item["fileEqual"].(bool)
		}
	}

	log.Println("[trace] service.signature.response.unmarshal.ok")
	// Mensaje según el estado HTTP
	message := "Verificación de firma completada correctamente."
	if resp.StatusCode != 200 {
		message = "Error en verificación de firma"
	}

	// Armar respuesta
	verificacion := map[string]interface{}{
		"statusCode": resp.StatusCode,
		"fileEqual":  fileEqual,
		"Message":    message,
	}

	log.Printf("[trace] service.signature.end | fileEqual=%v status=%d\n", fileEqual, resp.StatusCode)
	success := resp.StatusCode == 200
	return requestresponse.APIResponseDTO(success, resp.StatusCode, map[string]interface{}{"Verificacion": verificacion}, message)
}


/*package services

import (
	"os"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/verificacion_firma_digital_mid/models"
)

func PostVerificarFirma(archivo models.EmailAttachment, token string) requestresponse.APIResponse {
	// Establecer token como variable de entorno para el paquete request
	os.Setenv("Authorization", token)
	fmt.Println("Token establecido:", token)

	url := beego.AppConfig.String("FirmaElectronicaService") + "verify"

	payload := []models.PayloadVerificacion{
		{
			FileUp:    archivo.PdfBase64,
			Firma:     archivo.Firma,
			UrlFileUp: archivo.UrlFileUp,
		},
	}

	var resultado map[string]interface{}

	err := request.SendJson(url, "POST", &resultado, payload)
	if err != nil {
		beego.Error("Error en POST a verificador firma electrónica:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error en petición POST: %v", err))
	}

	return requestresponse.APIResponseDTO(true, 200, resultado, "Verificación de firma completada correctamente.")
}*/