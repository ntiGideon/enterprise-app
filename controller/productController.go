package controller

import (
	"Enterprise/helpers"
	"Enterprise/model"
	"Enterprise/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type ProductController struct {
	ProductService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (controller ProductController) CreateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	productModel := model.ProductModel{}
	helpers.ReadRequestBody(r, &productModel)
	userId := r.Context().Value("userId").(int)
	productModel.UserId = userId
	webResponse := controller.ProductService.CreateProduct(r.Context(), &productModel)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	productModel := model.ProductModel{}
	helpers.ReadRequestBody(r, &productModel)
	userId := r.Context().Value("userId").(int)
	productId := params.ByName("productId")
	id, _ := strconv.Atoi(productId)
	productModel.ProductId = id
	productModel.UserId = userId
	webResponse := controller.ProductService.UpdateProduct(r.Context(), &productModel)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller ProductController) GetProductById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, _ := strconv.Atoi(productId)
	webResponse := controller.ProductService.GetProductById(r.Context(), id)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller ProductController) DeleteProductById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, _ := strconv.Atoi(productId)
	userId := r.Context().Value("userId").(int)
	webResponse := controller.ProductService.DeleteProductById(r.Context(), id, userId)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	webResponse := controller.ProductService.GetAllProducts(r.Context())
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller ProductController) UpdateProductStock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	productModel := model.ProductStock{}
	helpers.ReadRequestBody(r, &productModel)
	productId := params.ByName("productId")
	userId := r.Context().Value("userId").(int)
	id, _ := strconv.Atoi(productId)
	productModel.ProductId = id
	productModel.UserId = userId
	webResponse := controller.ProductService.UpdateProductStock(r.Context(), &productModel)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}
