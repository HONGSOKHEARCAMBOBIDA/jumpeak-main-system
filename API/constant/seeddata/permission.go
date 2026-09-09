package seeddata

import "mysql/model"

var Permissions = []model.Permission{
	{
		Name:        "add.user",
		DisplayName: "Add User",
	},
	{
		Name:        "view.user",
		DisplayName: "View User",
	},
	{
		Name:        "edit.user",
		DisplayName: "Edit User",
	},
	{
		Name:        "view.role.has.permission",
		DisplayName: "view.role.has.permission",
	},
	{
		Name:        "add.role.has.permission",
		DisplayName: "add.role.has.permission",
	},
	{
		Name:        "delete.role.has.permission",
		DisplayName: "delete.role.has.permission",
	},
	{
		Name:        "view.company",
		DisplayName: "view.company",
	},
	{
		Name:        "add.company",
		DisplayName: "add.company",
	},
	{
		Name:        "update.company",
		DisplayName: "update.company",
	},
	{
		Name:        "view.Branch",
		DisplayName: "view.Branch",
	},
	{
		Name:        "add.Branch",
		DisplayName: "add.Branch",
	},
	{
		Name:        "update.Branch",
		DisplayName: "update.Branch",
	},
	{
		Name:        "view.Customer",
		DisplayName: "view.Customer",
	},
	{
		Name:        "add.Customer",
		DisplayName: "add.Customer",
	},
	{
		Name:        "update.Customer",
		DisplayName: "update.Customer",
	},
	{
		Name:        "view.Product",
		DisplayName: "view.Product",
	},
	{
		Name:        "add.Product",
		DisplayName: "add.Product",
	},
	{
		Name:        "update.Product",
		DisplayName: "update.Product",
	},
}
