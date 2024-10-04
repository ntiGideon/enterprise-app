package main

import (
	"Enterprise/config"
	"Enterprise/controller"
	"Enterprise/helpers"
	"Enterprise/router"
	"Enterprise/service"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic("Error loading .env file")
	}

	fmt.Printf("Server starting on PORT %s \n", os.Getenv("PORT"))

	db, err := config.ConnectDB()
	if err != nil {
		helpers.PanicAllErrors(err)
	}
	redisClient, err := config.ConnectRedis()

	defer db.Prisma.Disconnect()

	userService := service.NewUserService(db, redisClient)
	userController := controller.NewUserController(userService)
	categoryService := service.NewCategoryService(db)
	categoryController := controller.NewCategoryController(categoryService)
	productService := service.NewProductService(db)
	productController := controller.NewProductController(productService)

	routes := router.NewRouter(userController, categoryController, productController)

	server := http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverError := server.ListenAndServe()
	if serverError != nil {
		helpers.PanicAllErrors(serverError)
	}
}
