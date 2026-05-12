package services

import (
	"encoding/json"
	"log"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

type LambdaResponse struct {
	Status    string `json:"status"`     // "clean" o "infected"
	RawOutput string `json:"raw_output"` // salida de ClamAV
}

type LambdaRawResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"` // es un string que contiene JSON
}

func VerificarPDFBase64(pdfBase64 string) requestresponse.APIResponse {
	log.Println("[trace] service.clamav.start")
	if pdfBase64 == "" {
		return requestresponse.APIResponseDTO(false, 400, nil, "El contenido PDF base64 está vacío")
	}

	payload := map[string]string{"pdf_base64": pdfBase64}
	var rawResponse LambdaRawResponse

	//DOCKER CON CLAMDSCAN
	/*log.Println("[trace] service.clamav.request.send.start")
	err := request.SendJson(
		"http://localhost:8080/v1/verificar",
		"POST", &rawResponse, payload,
	)*/
	/*err := request.SendJson(
		"https://pruebasarchivo.portaloas.udistrital.edu.co/v1/verificar",
		"POST", &rawResponse, payload,
	)*/
	urlEscanear := beego.AppConfig.String("EscanearArchivo") + "verificar"
	log.Println("[trace] service.clamav.request.send.start")
	err := request.SendJson(
		urlEscanear,
		"POST", &rawResponse, payload,
	)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al llamar a Lambda: "+err.Error())
	}
	log.Printf("[trace] service.clamav.request.send.end | statusCode=%d\n", rawResponse.StatusCode)
	var lambdaResult LambdaResponse
	log.Println("[trace] service.clamav.response.unmarshal.start")
	if err := json.Unmarshal([]byte(rawResponse.Body), &lambdaResult); err != nil {
		log.Println("[trace] service.clamav.response.unmarshal.fail")
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al decodificar cuerpo de respuesta del antivirus: "+err.Error())
	}

	log.Printf("[trace] service.clamav.response.unmarshal.ok | result=%s\n", lambdaResult.Status)
	if lambdaResult.Status != "clean" && lambdaResult.Status != "infected" {
		return requestresponse.APIResponseDTO(false, 500, nil, "Respuesta inválida del antivirus")
	}

	return requestresponse.APIResponseDTO(true, 200, map[string]interface{}{
		"Virus": map[string]interface{}{
			"message":    "Verificación de virus completada correctamente.",
			"archive":    lambdaResult.Status,
			"statusCode": 200,
		},
	}, nil)
}
