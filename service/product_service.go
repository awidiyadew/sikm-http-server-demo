package service

import (
	"demo-app/model"
	"time"
)

func GetProduct() []model.Product {
	time.Sleep(1 * time.Second)
	return []model.Product{
		{
			Name: "iPhone X",
		},
		{
			Name: "iPhone 14",
		},
	}
}

func GetAdsProduct() []model.Product {
	time.Sleep(2 * time.Second)
	return []model.Product{
		{
			Name: "Google Pixel",
		},
		{
			Name: "Nokia",
		},
	}
}
