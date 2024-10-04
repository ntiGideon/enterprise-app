package service

import (
	"Enterprise/data"
	"Enterprise/helpers"
	"Enterprise/model"
	"Enterprise/prisma/db"
	"Enterprise/repository"
	"golang.org/x/net/context"
	"net/http"
)

type CategoryService struct {
	Db *db.PrismaClient
}

func NewCategoryService(db *db.PrismaClient) *CategoryService {
	return &CategoryService{Db: db}
}

func (p *CategoryService) CreateCategory(ctx context.Context, categoryModel *model.CategoryModel) *data.WebResponse {
	validator := helpers.RequestValidators(categoryModel)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}

	existingCategory, err := repository.ExistingCategoryByName(ctx, p.Db, categoryModel.Name)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	if existingCategory {
		return &data.WebResponse{
			Code:    http.StatusConflict,
			Message: "Category Name already exists",
			Data:    nil,
		}
	}

	_, err = p.Db.Category.CreateOne(
		db.Category.Name.Set(categoryModel.Name),
		db.Category.User.Link(db.User.ID.Equals(categoryModel.UserId)),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	err = repository.AuditLogs(ctx, p.Db, categoryModel.UserId, "Category created", "This action was performed by ")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &data.WebResponse{
		Code:    http.StatusCreated,
		Message: "Category Created",
		Data:    nil,
	}
}

func (p *CategoryService) GetCategoryById(ctx context.Context, categoryId int) *data.WebResponse {
	category, _ := p.Db.Category.FindUnique(db.Category.ID.Equals(categoryId)).Select(
		db.Category.ID.Field(),
		db.Category.Name.Field(),
	).Exec(ctx)

	if category == nil {
		return &data.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Category Not Found",
			Data:    nil,
		}
	}
	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Category found",
		Data: struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}{
			Id:   category.ID,
			Name: category.Name,
		},
	}
}

func (p *CategoryService) GetAllCategories(ctx context.Context) *data.WebResponse {
	categories, err := p.Db.Category.FindMany().Select(
		db.Category.ID.Field(),
		db.Category.Name.Field(),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	var CategoryResponses []model.CategoryResponse
	for _, category := range categories {
		CategoryResponses = append(CategoryResponses, model.CategoryResponse{
			Id:   category.ID,
			Name: category.Name,
		})
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Categories",
		Data:    CategoryResponses,
	}
}

func (p *CategoryService) UpdateCategory(ctx context.Context, category *model.CategoryModel) *data.WebResponse {
	validator := helpers.RequestValidators(category)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}
	existingCategory, err := repository.ExistingCategoryByName(ctx, p.Db, category.Name)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusForbidden,
			Message: err.Error(),
			Data:    nil,
		}
	}
	if existingCategory {
		return &data.WebResponse{
			Code:    http.StatusConflict,
			Message: "Category Name already exists",
			Data:    nil,
		}
	}
	_, err = p.Db.Category.FindUnique(db.Category.ID.Equals(category.CategoryId)).Update(db.Category.Name.Set(category.Name)).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	err = repository.AuditLogs(ctx, p.Db, category.UserId, "Category updated", "This action was performed by ")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Category Updated",
		Data:    nil,
	}

}

func (p *CategoryService) DeleteCategory(ctx context.Context, categoryId int, auditId int) *data.WebResponse {
	existingCategory, _ := p.Db.Category.FindUnique(db.Category.ID.Equals(categoryId)).Exec(ctx)
	if existingCategory == nil {
		return &data.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Category Not Found",
			Data:    nil,
		}
	}
	_, err := p.Db.Category.FindUnique(db.Category.ID.Equals(categoryId)).Delete().Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	err = repository.AuditLogs(ctx, p.Db, auditId, "Category deleted", "This action was performed by ")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Category Deleted",
		Data:    nil,
	}
}

func (p *CategoryService) AuditLogs(ctx context.Context) *data.WebResponse {
	logs, err := p.Db.AuditLog.FindMany().Select(
		db.AuditLog.Action.Field(),
		db.AuditLog.Details.Field(),
		db.AuditLog.CreatedAt.Field(),
	).With(db.AuditLog.User.Fetch().Select(db.User.Email.Field(), db.User.FirstName.Field())).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	var AuditLogResponses []model.AuditLogResponse
	for _, log := range logs {
		details, _ := log.Details()
		AuditLogResponses = append(AuditLogResponses, model.AuditLogResponse{
			Action:    log.Action,
			Details:   details,
			CreatedAt: log.CreatedAt,
			User: struct {
				Email     string `json:"email"`
				FirstName string `json:"firstName"`
			}{
				Email:     log.User().Email,
				FirstName: log.User().FirstName,
			},
		})
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "AuditLogs",
		Data:    AuditLogResponses,
	}

}
