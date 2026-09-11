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
	branchcontroller := controller.NewBranchController()
	customercontroller := controller.NewCustomerController()
	productcontroller := controller.NewProductController()
	invoicecontroller := controller.NewInvoiceController()
	paymentcontroller := controller.NewPaymentController()
	deptadjustmentcontroller := controller.NewDebtAdjustmentController()
	refundcontroller := controller.NewRefundController()
	customerledgercontroller := controller.NewCustomerledgerController()
	reportcontroller := controller.NewReportController()
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
		auth.GET(route.ViewCompanyNoPagination, middleware.PermissionMiddleware(permission.ViewCompany), companycontroller.GetCompanyNoPagination)

		// Branch
		auth.POST(route.AddBranch, middleware.PermissionMiddleware(permission.AddBranch), branchcontroller.Create)
		auth.PUT(route.UpdateBranch, middleware.PermissionMiddleware(permission.UpdateBranch), branchcontroller.Update)
		auth.GET(route.ViewBranchNoPagination, middleware.PermissionMiddleware(permission.ViewBranch), branchcontroller.GetBranchNoPagination)

		// User
		auth.POST(route.AddUser, middleware.PermissionMiddleware(permission.AddUser), authcontroller.Create)
		auth.PUT(route.EditUser, middleware.PermissionMiddleware(permission.EditUser), authcontroller.Update)

		// Customer
		auth.POST(route.AddCustomer, middleware.PermissionMiddleware(permission.AddCustomer), customercontroller.Create)
		auth.GET(route.ViewCustomer, middleware.PermissionMiddleware(permission.ViewCustomer), customercontroller.Get)
		auth.PUT(route.UpdateCustomer, middleware.PermissionMiddleware(permission.UpdateCustomer), customercontroller.Update)

		// Product
		auth.POST(route.AddProduct, middleware.PermissionMiddleware(permission.AddProduct), productcontroller.Create)
		auth.GET(route.ViewProduct, middleware.PermissionMiddleware(permission.ViewProduct), productcontroller.Get)
		auth.PUT(route.UpdateProduct, middleware.PermissionMiddleware(permission.UpdateProduct), productcontroller.Update)

		// Invoice
		auth.GET(route.ViewInvoice, middleware.PermissionMiddleware(permission.ViewInvoice), invoicecontroller.Get)
		auth.POST(route.AddInvoice, middleware.PermissionMiddleware(permission.AddInvoice), invoicecontroller.Create)
		auth.PUT(route.CancelInvoice, middleware.PermissionMiddleware(permission.CancelInvoice), invoicecontroller.Cancel)

		// Payment
		auth.GET(route.ViewPayment, middleware.PermissionMiddleware(permission.ViewPayment), paymentcontroller.Get)
		auth.POST(route.AddPayment, middleware.PermissionMiddleware(permission.AddPayment), paymentcontroller.Create)
		auth.PUT(route.VoidPayment, middleware.PermissionMiddleware(permission.VoidPayment), paymentcontroller.Void)

		// Refund
		auth.GET(route.ViewRefund, middleware.PermissionMiddleware(permission.ViewRefund), refundcontroller.Get)
		auth.POST(route.AddRefund, middleware.PermissionMiddleware(permission.AddRefund), refundcontroller.Create)

		// Debadjustment
		auth.GET(route.ViewDebAdjustment, middleware.PermissionMiddleware(permission.ViewDebAdjustment), deptadjustmentcontroller.Get)
		auth.POST(route.AddDebAdjustment, middleware.PermissionMiddleware(permission.AddDebAdjustment), deptadjustmentcontroller.Create)

		// CustomerLedger
		auth.GET(route.ViewCustomerledger, middleware.PermissionMiddleware(permission.ViewCustomerledger), customerledgercontroller.Get)

		// Reportcontroller
		auth.GET(route.ViewCustomerOutstadingReport, middleware.PermissionMiddleware(permission.ViewReport), reportcontroller.GetCustomerOutstandingReport)
	}
}
