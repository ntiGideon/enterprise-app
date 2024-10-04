package model

import "time"

type ProductModel struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
	Stock       int     `json:"stock" validate:"required"`
	CategoryId  []int   `json:"categoryId"`
	UserId      int     `json:"userId"`
	ProductId   int     `json:"productId"`
}

type ProductResponse struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProductStock struct {
	Stock     int `json:"stock" validate:"required"`
	ProductId int `json:"productId"`
	UserId    int `json:"userId"`
}
