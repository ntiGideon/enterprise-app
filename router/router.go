package router

import (
	"Enterprise/controller"
	"Enterprise/middleware"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(
	userController *controller.UserController,
	categoryController *controller.CategoryController,
	productController *controller.ProductController,
) *httprouter.Router {
	router := httprouter.New()

	var allowedRoles = []string{"ADMIN", "MANAGER", "EMPLOYEE"}
	var allowedRolesForAdmins = []string{"ADMIN"}
	var allowedRolesForManagers = []string{"ADMIN", "MANAGER"}

	router.POST("/api/admin/users/roles", middleware.RoleBasedAuthMiddleware(allowedRolesForAdmins, userController.RoleCreation))
	router.POST("/api/admin/users/create", middleware.RoleBasedAuthMiddleware(allowedRoles, userController.CreateUserByAdmin))
	router.PUT("/api/admin/users/update-info/:userId", middleware.RoleBasedAuthMiddleware(allowedRoles, userController.UpdateUserInfo))
	router.PUT("/api/admin/users/deactivate/:userId", middleware.RoleBasedAuthMiddleware(allowedRolesForAdmins, userController.DeactivateUser))
	router.PUT("/api/admin/users/delete/:userId", middleware.RoleBasedAuthMiddleware(allowedRoles, userController.DeleteUser))
	router.POST("/api/users/password", userController.CreateUserPassword)
	router.POST("/api/users/login", userController.Login)
	router.PUT("/api/users/change-info", middleware.RoleBasedAuthMiddleware(allowedRoles, userController.ChangeUserInfo))
	router.GET("/api/admin/users", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, userController.GetAllUsers))
	// AuditLogs
	router.GET("/api/admin/logs", middleware.RoleBasedAuthMiddleware(allowedRolesForAdmins, categoryController.AuditLogs))

	// Categories
	router.POST("/api/category/create", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, categoryController.CreateCategory))
	router.GET("/api/category/:categoryId", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, categoryController.GetCategoryById))
	router.GET("/api/category", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, categoryController.GetAllCategories))
	router.PUT("/api/category/update/:categoryId", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, categoryController.UpdateCategory))
	router.DELETE("/api/category/delete/:categoryId", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, categoryController.DeleteCategory))

	// Product
	router.POST("/api/product/create", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, productController.CreateProduct))
	router.PUT("/api/product/update/:productId", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, productController.UpdateProduct))
	router.GET("/api/product/:productId", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, productController.GetProductById))
	router.DELETE("/api/product/delete/:productId", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, productController.DeleteProductById))
	router.GET("/api/product", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, productController.GetAllProducts))
	router.PUT("/api/product-stock/:productId", middleware.RoleBasedAuthMiddleware(allowedRolesForManagers, productController.UpdateProductStock))

	return router
}
