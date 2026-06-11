package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/verificacion_firma_digital_mid/models"
)

func PostVerificarFirma(archivo models.EmailAttachment, token string) requestresponse.APIResponse {
	url := beego.AppConfig.String("FirmaElectronicaService") + "verify"

	payload := []models.PayloadVerificacion{
		{
			FileUp:    archivo.PdfBase64,
			Firma:     archivo.Firma,
			UrlFileUp: archivo.UrlFileUp,
		},
	}

	// Serializar el payload a JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		beego.Error("Error al serializar JSON:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al serializar JSON: "+err.Error())
	}

	// Crear la solicitud HTTP con body y header Authorization
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		beego.Error("Error al crear solicitud HTTP:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al crear solicitud HTTP: "+err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	//req.Header.Add("Authorization", "Bearer 299a6897-f955-39c1-b4ce-9fd731387b3d") // Usar Add para evitar reemplazar si ya existe

	client := &http.Client{}
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
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var resultado map[string]interface{}
	if err := json.Unmarshal(body, &resultado); err != nil {
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

	success := resp.StatusCode == 200
	return requestresponse.APIResponseDTO(success, resp.StatusCode, map[string]interface{}{"Verificacion": verificacion}, message)
}
