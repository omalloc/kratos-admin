package data

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/omalloc/kratos-admin/internal/biz"
)

// SeedData 初始化基础数据
func SeedData(db *gorm.DB, logger log.Logger) error {
	log := log.NewHelper(logger)
	log.Info("开始初始化基础数据...")

	// 检查是否已经初始化过
	var count int64
	if err := db.Model(&biz.User{}).Count(&count).Error; err != nil {
		return err
	}

	// 如果已经有数据，跳过初始化
	if count > 0 {
		log.Info("数据库中已有数据，跳过初始化")
		return nil
	}

	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 初始化权限
		if err := seedPermissions(tx, log); err != nil {
			return err
		}

		// 2. 初始化角色
		if err := seedRoles(tx, log); err != nil {
			return err
		}

		// 3. 初始化角色权限关联
		if err := seedRolePermissions(tx, log); err != nil {
			return err
		}

		// 4. 初始化管理员用户
		if err := seedAdminUser(tx, log); err != nil {
			return err
		}

		// 5. 初始化菜单
		if err := seedMenus(tx, log); err != nil {
			return err
		}

		log.Info("基础数据初始化完成")
		return nil
	})
}

// seedPermissions 初始化权限数据
func seedPermissions(tx *gorm.DB, log *log.Helper) error {
	permissions := []*biz.Permission{
		{
			ID:       1,
			Name:     "user",
			Alias:    "用户管理",
			Describe: "用户管理权限",
			Status:   1,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建用户", Checked: true},
				{Key: "READ", Describe: "查看用户", Checked: true},
				{Key: "UPDATE", Describe: "更新用户", Checked: true},
				{Key: "DELETE", Describe: "删除用户", Checked: true},
			},
		},
		{
			ID:       2,
			Name:     "role",
			Alias:    "角色管理",
			Describe: "角色管理权限",
			Status:   1,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建角色", Checked: true},
				{Key: "READ", Describe: "查看角色", Checked: true},
				{Key: "UPDATE", Describe: "更新角色", Checked: true},
				{Key: "DELETE", Describe: "删除角色", Checked: true},
			},
		},
		{
			ID:       3,
			Name:     "permission",
			Alias:    "权限管理",
			Describe: "权限管理权限",
			Status:   1,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建权限", Checked: true},
				{Key: "READ", Describe: "查看权限", Checked: true},
				{Key: "UPDATE", Describe: "更新权限", Checked: true},
				{Key: "DELETE", Describe: "删除权限", Checked: true},
			},
		},
		{
			ID:       4,
			Name:     "menu",
			Alias:    "菜单管理",
			Describe: "菜单管理权限",
			Status:   1,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建菜单", Checked: true},
				{Key: "READ", Describe: "查看菜单", Checked: true},
				{Key: "UPDATE", Describe: "更新菜单", Checked: true},
				{Key: "DELETE", Describe: "删除菜单", Checked: true},
			},
		},
		{
			ID:       5,
			Name:     "system",
			Alias:    "系统管理",
			Describe: "系统管理权限",
			Status:   1,
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看系统", Checked: true},
				{Key: "UPDATE", Describe: "更新系统", Checked: true},
			},
		},
	}

	for _, perm := range permissions {
		if err := tx.Create(perm).Error; err != nil {
			log.Errorf("创建权限失败: %v", err)
			return err
		}
	}

	log.Info("权限数据初始化完成")
	return nil
}

// seedRoles 初始化角色数据
func seedRoles(tx *gorm.DB, log *log.Helper) error {
	roles := []*biz.Role{
		{
			ID:       1,
			Name:     "super_admin",
			Alias:    "超级管理员",
			Describe: "系统超级管理员，拥有所有权限",
			Status:   1,
		},
		{
			ID:       2,
			Name:     "admin",
			Alias:    "管理员",
			Describe: "系统管理员，拥有大部分权限",
			Status:   1,
		},
		{
			ID:       3,
			Name:     "user",
			Alias:    "普通用户",
			Describe: "普通用户，拥有基本权限",
			Status:   1,
		},
	}

	for _, role := range roles {
		if err := tx.Create(role).Error; err != nil {
			log.Errorf("创建角色失败: %v", err)
			return err
		}
	}

	log.Info("角色数据初始化完成")
	return nil
}

// seedRolePermissions 初始化角色权限关联
func seedRolePermissions(tx *gorm.DB, log *log.Helper) error {
	// 超级管理员拥有所有权限
	superAdminPermissions := []*biz.RolePermission{
		{
			RoleID: 1,
			PermID: 1,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建用户", Checked: true},
				{Key: "READ", Describe: "查看用户", Checked: true},
				{Key: "UPDATE", Describe: "更新用户", Checked: true},
				{Key: "DELETE", Describe: "删除用户", Checked: true},
			},
		},
		{
			RoleID: 1,
			PermID: 2,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建角色", Checked: true},
				{Key: "READ", Describe: "查看角色", Checked: true},
				{Key: "UPDATE", Describe: "更新角色", Checked: true},
				{Key: "DELETE", Describe: "删除角色", Checked: true},
			},
		},
		{
			RoleID: 1,
			PermID: 3,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建权限", Checked: true},
				{Key: "READ", Describe: "查看权限", Checked: true},
				{Key: "UPDATE", Describe: "更新权限", Checked: true},
				{Key: "DELETE", Describe: "删除权限", Checked: true},
			},
		},
		{
			RoleID: 1,
			PermID: 4,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建菜单", Checked: true},
				{Key: "READ", Describe: "查看菜单", Checked: true},
				{Key: "UPDATE", Describe: "更新菜单", Checked: true},
				{Key: "DELETE", Describe: "删除菜单", Checked: true},
			},
		},
		{
			RoleID: 1,
			PermID: 5,
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看系统", Checked: true},
				{Key: "UPDATE", Describe: "更新系统", Checked: true},
			},
		},
	}

	// 管理员拥有大部分权限
	adminPermissions := []*biz.RolePermission{
		{
			RoleID: 2,
			PermID: 1,
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建用户", Checked: true},
				{Key: "READ", Describe: "查看用户", Checked: true},
				{Key: "UPDATE", Describe: "更新用户", Checked: true},
			},
		},
		{
			RoleID: 2,
			PermID: 2,
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看角色", Checked: true},
				{Key: "UPDATE", Describe: "更新角色", Checked: true},
			},
		},
		{
			RoleID: 2,
			PermID: 4,
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看菜单", Checked: true},
			},
		},
		{
			RoleID: 2,
			PermID: 5,
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看系统", Checked: true},
			},
		},
	}

	// 普通用户只有基本权限
	userPermissions := []*biz.RolePermission{
		{
			RoleID: 3,
			PermID: 1,
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看用户", Checked: true},
			},
		},
		{
			RoleID: 3,
			PermID: 4,
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看菜单", Checked: true},
			},
		},
	}

	allPermissions := append(superAdminPermissions, adminPermissions...)
	allPermissions = append(allPermissions, userPermissions...)

	for _, rolePerm := range allPermissions {
		if err := tx.Create(rolePerm).Error; err != nil {
			log.Errorf("创建角色权限关联失败: %v", err)
			return err
		}
	}

	log.Info("角色权限关联初始化完成")
	return nil
}

// seedAdminUser 初始化管理员用户
func seedAdminUser(tx *gorm.DB, log *log.Helper) error {
	// 创建超级管理员用户
	adminUser := &biz.User{
		ID:       1,
		Username: "admin",
		Password: "$2a$10$aEPHKnracWQOU8wQwA2SsuIeLyae/qlcQDp75j3TPQlPuVUkCtDwa", // 默认密码: 12346578
		Email:    "admin@example.com",
		Nickname: "超级管理员",
		Bio:      "系统超级管理员",
		Status:   1,
	}

	if err := tx.Create(adminUser).Error; err != nil {
		log.Errorf("创建管理员用户失败: %v", err)
		return err
	}

	// 为管理员用户分配超级管理员角色
	userRole := &biz.UserRole{
		UserID:    1,
		RoleID:    1, // 超级管理员角色
		CreatedAt: time.Now(),
	}

	if err := tx.Create(userRole).Error; err != nil {
		log.Errorf("分配管理员角色失败: %v", err)
		return err
	}

	log.Info("管理员用户初始化完成")
	return nil
}

// seedMenus 初始化菜单数据
func seedMenus(tx *gorm.DB, log *log.Helper) error {
	menus := []*biz.Menu{
		{
			ID:           1,
			PID:          0,
			PermissionID: 0,
			Name:         "仪表盘",
			Icon:         "icon-app-box",
			Path:         "/",
			SortBy:       1,
			Hidden:       false,
			Status:       1,
		},
		{
			ID:           2,
			PID:          0,
			PermissionID: 0,
			Name:         "系统管理",
			Icon:         "icon-config",
			Path:         "/admin",
			SortBy:       2,
			Hidden:       false,
			Status:       1,
		},
		{
			ID:           3,
			PID:          2,
			PermissionID: 1,
			Name:         "用户",
			Icon:         "icon-user",
			Path:         "/admin/user",
			SortBy:       1,
			Hidden:       false,
			Status:       1,
		},
		{
			ID:           4,
			PID:          2,
			PermissionID: 2,
			Name:         "角色",
			Icon:         "icon-addteam",
			Path:         "/admin/role",
			SortBy:       2,
			Hidden:       false,
			Status:       1,
		},
		{
			ID:           5,
			PID:          2,
			PermissionID: 3,
			Name:         "权限",
			Icon:         "icon-securityscan",
			Path:         "/admin/permission",
			SortBy:       3,
			Hidden:       false,
			Status:       1,
		},
		{
			ID:           6,
			PID:          2,
			PermissionID: 4,
			Name:         "菜单",
			Icon:         "icon-resource",
			Path:         "/admin/menu",
			SortBy:       4,
			Hidden:       false,
			Status:       1,
		},
	}

	for _, menu := range menus {
		if err := tx.Create(menu).Error; err != nil {
			log.Errorf("创建菜单失败: %v", err)
			return err
		}
	}

	log.Info("菜单数据初始化完成")
	return nil
}
