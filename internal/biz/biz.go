package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	// rbac modules.
	NewUserUsecase,
	NewRoleUsecase,
	NewPermissionUsecase,
	NewMenuUsecase,
	NewCrontabUsecase,
)
