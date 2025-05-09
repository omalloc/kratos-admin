package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewConsoleService,
	// User, Role, Permissions,
	NewUserService,
	NewRoleService,
	NewPermissionService,
	NewPassportService,
)
