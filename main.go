package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"demo-app/model"
	"demo-app/service"
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
		if r.Method != "POST" {
			http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

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

func main() {
	http.HandleFunc("/hello", GetHello())
	http.HandleFunc("/hello/html", GetHelloHTML())
	http.HandleFunc("/contact/list", GetContacts())
	http.HandleFunc("/contact/add", PostContacts())

	http.ListenAndServe(":3000", nil)
}
