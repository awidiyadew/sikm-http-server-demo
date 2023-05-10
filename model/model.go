package model

type Contact struct {
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	ImgURL string `json:"img_url"`
}

type ErrorResp struct {
	Message string `json:"message"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	Name string `json:"name"`
}
