package model

import "time"

type CategoryModel struct {
	Name       string `json:"name" validate:"required,min=3,max=50"`
	UserId     int    `json:"userId"`
	CategoryId int    `json:"categoryId"`
}

type CategoryResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
