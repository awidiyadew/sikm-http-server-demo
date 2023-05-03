package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func x() {
	file, err := os.OpenFile("data/contact.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("error open file: ", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("\n08512212319_Jane_https://img.com/jane")
	if err != nil {
		fmt.Println("failed to add contact: ", err)
		return
	}

	contacts, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file: ", err)
		return
	}

	str := string(contacts)
	arr := strings.Split(str, "\n")

	fmt.Println(len(arr), arr, str)
}
