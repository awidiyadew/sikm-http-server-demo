package model

type Contact struct {
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	ImgURL string `json:"img_url"`
}

type ErrorResp struct {
	message string `json:"message"`
}
