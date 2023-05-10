package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"demo-app/middleware"
	"demo-app/model"
	"demo-app/service"
)

// anggap ini data dari DB
const (
	username = "johndoe"
	password = "john123"
)

func GetHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		text := "Hello"
		if name != "" {
			text += " " + name
		} else {
			text += " " + "World"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(text))
	}
}

func GetHelloHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		tmpl, err := template.ParseFiles("template/hello.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// bisa pakai map atau pake struct
		data := map[string]string{
			"name": name,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetContacts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contacts, _ := service.GetContacts()
		jsonResp, err := json.Marshal(contacts)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		// order is matters, set content type dlu baru write status code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

func PostContacts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// membaca data username yg disimpan oleh middleware Auth
		loggedInUser := r.Context().Value("username")
		fmt.Println("new contact added by ", loggedInUser)

		reqBody, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		var c model.Contact
		err = json.Unmarshal(reqBody, &c)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		insertedContact, err := service.AddContact(c)
		if err != nil {
			// http.Error(w, "failed add contact", http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			errBody, _ := json.Marshal(map[string]string{
				"message": err.Error(),
			})
			w.Write(errBody)
			return
		}
		resp, _ := json.Marshal(insertedContact)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func HandleLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload model.LoginData
		json.NewDecoder(r.Body).Decode(&payload)

		// TODO: check ke DB apakah username dan pwd sesuai?

		if payload.Username != username || payload.Password != password {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest) // bisa 400 atau 401
			json.NewEncoder(w).Encode(model.ErrorResp{
				Message: "invalid username or password",
			})
			return
		}

		// username dan pwd sudah benar
		cookieUsername := http.Cookie{
			Name:   "username",
			Value:  payload.Username,
			MaxAge: 3600, // 3600s atau 1 jam
		}

		cookieData := http.Cookie{
			Name:   "data",
			Value:  "hello world",
			MaxAge: 3600, // 3600s atau 1 jam
		}

		http.SetCookie(w, &cookieUsername)
		http.SetCookie(w, &cookieData)

		w.Write([]byte(payload.Username + " successfully logged in"))
	})
}

func HandleLogout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: check is cookie exist

		c := http.Cookie{
			Name:   "username",
			MaxAge: -1,
		}
		http.SetCookie(w, &c)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("logout success"))
	})
}

func HandleGetProducts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch := make(chan []model.Product)
		go func() {
			ads := service.GetAdsProduct()
			ch <- ads
		}()

		go func() {
			products := service.GetProduct()
			ch <- products
		}()

		allProducts := []model.Product{}
		for i := 0; i < 2; i++ {
			prd := <-ch
			allProducts = append(allProducts, prd...)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(allProducts)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", middleware.ValidateMethod(http.MethodGet, GetHello()))
	mux.Handle("/hello/html", middleware.ValidateMethod(http.MethodGet, GetHelloHTML()))
	mux.Handle("/contact/list", middleware.ValidateMethod(http.MethodGet, GetContacts()))
	mux.Handle("/contact/add", middleware.ValidateMethod(http.MethodPost, middleware.Auth(PostContacts())))

	mux.Handle("/login", middleware.ValidateMethod(http.MethodPost, HandleLogin()))
	mux.Handle("/logout", middleware.ValidateMethod(http.MethodPost, HandleLogout()))

	mux.Handle("/products", middleware.ValidateMethod(http.MethodGet, HandleGetProducts()))

	http.ListenAndServe(":3000", mux)
}
