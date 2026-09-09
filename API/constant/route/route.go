package route

const (

	// Authentication
	Login     = "login"
	LoginByQr = "loginbyqr"
	Refresh   = "refresh"
	Logout    = "logout"

	ViewUserData   = "view.user.data"
	AddUser        = "add.user"
	ViewUser       = "view.user"
	EditUser       = "edit.user/:id"
	ChangePassword = "change.password"

	// RoleHasPermission
	ViewRoleHasPermission   = "view.role.has.permission/:id"
	AddRoleHasPermission    = "add.role.has.permission"
	DeleteRoleHasPermission = "delete.role.has.permission"

	// Role
	ViewRole = "view.role"
	EditRole = "edit.role/:id"

	// Company
	ViewCompany             = "view.company"
	ViewCompanyNoPagination = "view.company.no.pagination"
	AddCompany              = "add.company"
	UpdateCompany           = "update.company/:id"

	// Branch
	ViewBranch             = "view.Branch"
	ViewBranchNoPagination = "view.Branch.no.pagination/:id"
	AddBranch              = "add.Branch"
	UpdateBranch           = "update.Branch/:id"

	// Customer
	ViewCustomer   = "view.Customer"
	AddCustomer    = "add.Customer"
	UpdateCustomer = "update.Customer/:id"
)
