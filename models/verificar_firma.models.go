package models

/*type EmailAttachment struct {
	Hash      string `json:"hash"`
	PdfBase64 string `json:"pdf_base64"`
}*/

type EmailAttachment struct {
	PdfBase64 string `json:"pdf_base64"` // equivale a fileUp
	Firma     string `json:"firma"`      // hash de firma electrónica
	UrlFileUp string `json:"urlFileUp"`  // url del archivo original
}

type PayloadVerificacion struct {
	Firma     string `json:"firma"`
	FileUp    string `json:"fileUp"`
	UrlFileUp string `json:"urlFileUp"`
}
