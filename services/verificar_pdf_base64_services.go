package services

import (
	"bytes"
	"encoding/json"
	//"errors"
	"net/http"

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
	if pdfBase64 == "" {
		return requestresponse.APIResponseDTO(false, 400, nil, "El contenido PDF base64 está vacío")
	}

	payload := map[string]string{"pdf_base64": pdfBase64}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al serializar el payload: "+err.Error())
	}

	/*resp, err := http.Post("http://localhost:9000/2015-03-31/functions/function/invocations", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al llamar a Lambda: "+err.Error())
	}
	defer resp.Body.Close()

	var lambdaResult LambdaResponse
	if err := json.NewDecoder(resp.Body).Decode(&lambdaResult); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al decodificar respuesta del antivirus: "+err.Error())
	}*/

	resp, err := http.Post("http://localhost:9000/2015-03-31/functions/function/invocations", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al llamar a Lambda: "+err.Error())
	}
	defer resp.Body.Close()

	var rawResponse LambdaRawResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al decodificar respuesta bruta del antivirus: "+err.Error())
	}

	var lambdaResult LambdaResponse
	if err := json.Unmarshal([]byte(rawResponse.Body), &lambdaResult); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al decodificar cuerpo de respuesta del antivirus: "+err.Error())
	}


	// Verificación de estado
	if lambdaResult.Status != "clean" && lambdaResult.Status != "infected" {
		return requestresponse.APIResponseDTO(false, 500, nil, "Respuesta inválida del antivirus")
	}
	

	// Construir respuesta
	return requestresponse.APIResponseDTO(true, 200, map[string]interface{}{
		"status":     lambdaResult.Status,
		"raw_output": lambdaResult.RawOutput,
	}, nil)
}
