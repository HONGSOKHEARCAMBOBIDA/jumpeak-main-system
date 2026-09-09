package permission

const (
	// user
	AddUser        = "add.user"
	ViewUser       = "view.user"
	EditUser       = "edit.user"
	ChangePassword = "change.password"

	// RoleHasPermission
	ViewRoleHasPermission   = "view.role.has.permission"
	AddRoleHasPermission    = "add.role.has.permission"
	DeleteRoleHasPermission = "delete.role.has.permission"

	// Company
	ViewCompany   = "view.company"
	AddCompany    = "add.company"
	UpdateCompany = "update.company"

	// Branch
	ViewBranch   = "view.Branch"
	AddBranch    = "add.Branch"
	UpdateBranch = "update.Branch"

	// Customer
	ViewCustomer   = "view.Customer"
	AddCustomer    = "add.Customer"
	UpdateCustomer = "update.Customer"

	// Product
	ViewProduct   = "view.Product"
	AddProduct    = "add.Product"
	UpdateProduct = "update.Product"

	// Invoice
	ViewInvoice   = "view.Invoice"
	AddInvoice    = "add.Invoice"
	CancelInvoice = "Cancel.Invoice"

	// Payment
	ViewPayment = "view.Payment"
	AddPayment  = "add.Payment"
	VoidPayment = "Void.Payment"

	// DebAdjustment
	ViewDebAdjustment = "view.DebAdjustment"
	AddDebAdjustment  = "add.DebAdjustment"

	// Refund
	ViewRefund = "view.Refund"
	AddRefund  = "add.Refund"
)
