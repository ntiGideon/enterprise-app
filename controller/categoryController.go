package controller

import (
	"Enterprise/helpers"
	"Enterprise/model"
	"Enterprise/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type CategoryController struct {
	CategoryService *service.CategoryService
}

func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{
		CategoryService: categoryService,
	}
}

func (controller CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	categoryDto := model.CategoryModel{}
	helpers.ReadRequestBody(r, &categoryDto)
	userId := r.Context().Value("userId").(int)
	categoryDto.UserId = userId

	webResponse := controller.CategoryService.CreateCategory(r.Context(), &categoryDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller CategoryController) GetCategoryById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	id, _ := strconv.Atoi(categoryId)
	webResponse := controller.CategoryService.GetCategoryById(r.Context(), id)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	userId := r.Context().Value("userId").(int)
	id, _ := strconv.Atoi(categoryId)
	webResponse := controller.CategoryService.DeleteCategory(r.Context(), id, userId)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	webResponse := controller.CategoryService.GetAllCategories(r.Context())
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller CategoryController) AuditLogs(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	webResponse := controller.CategoryService.AuditLogs(r.Context())
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	categoryDto := model.CategoryModel{}
	helpers.ReadRequestBody(r, &categoryDto)
	userId := r.Context().Value("userId").(int)
	categoryDto.UserId = userId
	id, _ := strconv.Atoi(categoryId)
	categoryDto.CategoryId = id
	webResponse := controller.CategoryService.UpdateCategory(r.Context(), &categoryDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}
