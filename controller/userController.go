package controller

import (
	"Enterprise/helpers"
	"Enterprise/model"
	"Enterprise/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type UserController struct {
	UserService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (controller *UserController) RoleCreation(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	roleDto := model.RoleCreationModel{}
	helpers.ReadRequestBody(r, &roleDto)
	userId := r.Context().Value("userId").(int)
	roleDto.AuditId = userId

	webResponse := controller.UserService.RoleCreation(r.Context(), &roleDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) CreateUserByAdmin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := model.UserCreationModel{}
	helpers.ReadRequestBody(r, &userDto)
	userId := r.Context().Value("userId").(int)
	userDto.UserId = userId

	webResponse := controller.UserService.CreateUserByAdmin(r.Context(), &userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) CreateUserPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := model.UserPasswordCreationModel{}
	helpers.ReadRequestBody(r, &userDto)

	code := r.URL.Query().Get("code")
	userDto.Code = code

	webResponse := controller.UserService.CreateUserPassword(r.Context(), &userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) UpdateUserInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := model.UserCreationModel{}
	helpers.ReadRequestBody(r, &userDto)
	userId := params.ByName("userId")
	id, _ := strconv.Atoi(userId)
	userDto.UserId = id
	webResponse := controller.UserService.UpdateUserInfo(r.Context(), &userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) DeactivateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId, _ := strconv.Atoi(params.ByName("userId"))
	webResponse := controller.UserService.DeactivateUser(r.Context(), userId)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId, _ := strconv.Atoi(params.ByName("userId"))
	webResponse := controller.UserService.DeleteUser(r.Context(), userId)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := model.LoginUserModel{}
	helpers.ReadRequestBody(r, &userDto)
	webResponse := controller.UserService.Login(r.Context(), &userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) ChangeUserInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := model.UpdateUserInfoModel{}
	helpers.ReadRequestBody(r, &userDto)
	userId := r.Context().Value("userId").(int)
	userDto.UserId = userId
	webResponse := controller.UserService.ChangeUserInfo(r.Context(), &userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	webResponse := controller.UserService.GetAllUsers(r.Context())
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}
