package routes

import (
	"mysql/constant/permission"
	"mysql/constant/route"
	"mysql/controller"
	"mysql/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authcontroller := controller.NewAuthController()
	rolehaspermissioncontroller := controller.NewRoleHasPermissionController()
	companycontroller := controller.NewCompanyController()
	public := r.Group("/")
	public.Use(middleware.APIKeyAuth())
	{
		public.POST(route.Login, authcontroller.Login)
		public.POST(route.Refresh, authcontroller.Refresh)
	}
	auth := r.Group("/")
	auth.Use(middleware.APIKeyAuth())
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET(route.ViewRole, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetRole)
		auth.GET(route.ViewUserData, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetUserData)

		// RoleHasPermission
		auth.GET(route.ViewRoleHasPermission, middleware.PermissionMiddleware(permission.ViewRoleHasPermission), rolehaspermissioncontroller.GetRolePermission)
		auth.POST(route.AddRoleHasPermission, middleware.PermissionMiddleware(permission.AddRoleHasPermission), rolehaspermissioncontroller.CreateRoleHasPermission)
		auth.DELETE(route.DeleteRoleHasPermission, middleware.PermissionMiddleware(permission.DeleteRoleHasPermission), rolehaspermissioncontroller.DeleteRoleHasPermission)

		// Company
		auth.GET(route.ViewCompany, middleware.PermissionMiddleware(permission.ViewCompany), companycontroller.Get)
		auth.POST(route.AddCompany, middleware.PermissionMiddleware(permission.AddCompany), companycontroller.Create)
		auth.PUT(route.UpdateCompany, middleware.PermissionMiddleware(permission.UpdateCompany), companycontroller.Update)
	}
}
