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

type ProductService struct {
	Db *db.PrismaClient
}

func NewProductService(db *db.PrismaClient) *ProductService {
	return &ProductService{Db: db}
}

func (p *ProductService) CreateProduct(ctx context.Context, productDto *model.ProductModel) *data.WebResponse {
	validator := helpers.RequestValidators(productDto)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}

	existingProduct := repository.ExistingProductByName(ctx, p.Db, productDto.Name)
	if existingProduct {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Product already exists",
			Data:    nil,
		}
	}
	product, err := p.Db.Product.CreateOne(
		db.Product.Name.Set(productDto.Name),
		db.Product.Price.Set(float64(productDto.Price)),
		db.Product.Stock.Set(productDto.Stock),
		db.Product.Description.Set(productDto.Description),
		db.Product.User.Link(db.User.ID.Equals(productDto.UserId)),
		//db.Product.Categories.Link(db.Category.And(db.Category.ID.In(productDto.CategoryId))),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	for _, categoryId := range productDto.CategoryId {
		_, err = p.Db.ProductOnCategory.CreateOne(
			db.ProductOnCategory.Product.Link(db.Product.ID.Equals(product.ID)),
			db.ProductOnCategory.Category.Link(db.Category.ID.Equals(categoryId))).Exec(ctx)
	}

	err = repository.AuditLogs(ctx, p.Db, productDto.UserId, "Product Created", "This action was performed by ")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &data.WebResponse{
		Code:    http.StatusCreated,
		Message: "Product created",
		Data:    nil,
	}
}

func (p *ProductService) UpdateProduct(ctx context.Context, productDto *model.ProductModel) *data.WebResponse {
	validator := helpers.RequestValidators(productDto)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}
	existingProduct := repository.ExistingProductByName(ctx, p.Db, productDto.Name)
	if existingProduct {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Product already exists",
			Data:    nil,
		}
	}
	_, err := p.Db.Product.FindUnique(db.Product.ID.Equals(productDto.ProductId)).Update(
		db.Product.Stock.Set(productDto.Stock),
		db.Product.Name.Set(productDto.Name),
		db.Product.Description.Set(productDto.Description),
		db.Product.Price.Set(float64(productDto.Price)),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Product updated",
		Data:    nil,
	}
}

func (p *ProductService) GetProductById(ctx context.Context, productId int) *data.WebResponse {
	productExist, _ := p.Db.Product.FindUnique(db.Product.ID.Equals(productId)).Exec(ctx)
	if productExist == nil {
		return &data.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Product not found",
			Data:    nil,
		}
	}

	description, _ := productExist.Description()
	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Product found",
		Data: model.ProductResponse{
			Id:          productId,
			Name:        productExist.Name,
			Description: description,
			Price:       float32(productExist.Price),
			Stock:       productExist.Stock,
			CreatedAt:   productExist.CreatedAt,
		},
	}
}

func (p *ProductService) GetAllProducts(ctx context.Context) *data.WebResponse {
	products, _err := p.Db.Product.FindMany().Exec(ctx)
	if _err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: _err.Error(),
			Data:    nil,
		}
	}
	var ProductResponses []model.ProductResponse
	for _, product := range products {
		description, _ := product.Description()
		ProductResponses = append(ProductResponses, model.ProductResponse{
			Id:          product.ID,
			Name:        product.Name,
			Description: description,
			Price:       float32(product.Price),
			Stock:       product.Stock,
			CreatedAt:   product.CreatedAt,
		})
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Products found",
		Data:    ProductResponses,
	}
}

// DeleteProductById TODO fix cascading issue
func (p *ProductService) DeleteProductById(ctx context.Context, productId int, userId int) *data.WebResponse {
	productExist, _ := p.Db.Product.FindUnique(db.Product.ID.Equals(productId)).Exec(ctx)
	if productExist == nil {
		return &data.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Product not found",
			Data:    nil,
		}
	}

	_, err := p.Db.Product.FindUnique(db.Product.ID.Equals(productId)).Delete().Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	err = repository.AuditLogs(ctx, p.Db, userId, "Product Deleted", "This action was performed by ")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Product deleted",
		Data:    nil,
	}
}

func (p *ProductService) UpdateProductStock(ctx context.Context, productStock *model.ProductStock) *data.WebResponse {
	validator := helpers.RequestValidators(productStock)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}

	_, err := p.Db.Product.FindUnique(db.Product.ID.Equals(productStock.ProductId)).Update(
		db.Product.Stock.Set(productStock.Stock),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	err = repository.AuditLogs(ctx, p.Db, productStock.UserId, "Product Deleted", "This action was performed by ")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Product stock updated",
		Data:    nil,
	}
}
