package service

import (
	"demo-app/model"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetContacts() ([]model.Contact, error) {
	// contacts := []model.Contact{
	// 	{
	// 		Phone:  "085934567890",
	// 		Name:   "Dewa",
	// 		ImgURL: "https://cataas.com/cat/says/DEWA",
	// 	},
	// 	{
	// 		Phone:  "085934512345",
	// 		Name:   "John Doe",
	// 		ImgURL: "https://cataas.com/cat/says/JOHN DOE",
	// 	},
	// }

	// buka file dengan akses readonly
	file, err := os.OpenFile("data/contact.txt", os.O_RDONLY, 0644)
	if err != nil {
		return []model.Contact{}, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return []model.Contact{}, err
	}

	var contacts []model.Contact
	strs := strings.Split(string(content), "\n")
	for _, rowContact := range strs {
		data := strings.Split(rowContact, "_")
		contact := model.Contact{
			Phone:  data[0],
			Name:   data[1],
			ImgURL: data[2],
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func AddContact(c model.Contact) (model.Contact, error) {
	file, err := os.OpenFile("data/contact.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return model.Contact{}, err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "\n%v_%v_%v", c.Phone, c.Name, c.ImgURL)
	if err != nil {
		return model.Contact{}, err
	}

	return c, nil
}
