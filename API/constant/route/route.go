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
	EditUser       = "edit.user"
	ChangePassword = "change.password"

	// RoleHasPermission
	ViewRoleHasPermission   = "view.role.has.permission/:id"
	AddRoleHasPermission    = "add.role.has.permission"
	DeleteRoleHasPermission = "delete.role.has.permission"

	// Role
	ViewRole = "view.role"
	EditRole = "edit.role/:id"

	// Company
	ViewCompany   = "view.company"
	AddCompany    = "add.company"
	UpdateCompany = "update.company/:id"
)
