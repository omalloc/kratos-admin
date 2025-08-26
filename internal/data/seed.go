package data

import (
	"fmt"
	"slices"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/kratos-admin/api/console/administration"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/omalloc/kratos-admin/internal/biz"
)

type Seed struct {
	db    *gorm.DB
	idGen *snowflake.Node
}

// SeedData 初始化基础数据
func SeedData(db *gorm.DB) error {
	log.Debug("开始初始化基础数据...")

	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Errorf("sonyflake create failed: %v", err)
		return fmt.Errorf("sonyflake create failed: %v", err)
	}

	sd := &Seed{
		db:    db,
		idGen: node,
	}

	// 检查是否已经初始化过
	// 如果已经有数据，跳过初始化
	if skip, err := sd.Check(); err != nil {
		return err
	} else if skip {
		log.Debug("数据库中已有数据，跳过初始化")
		return nil
	}

	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 初始化权限
		permissions, err := sd.seedPermissions(tx)
		if err != nil {
			return err
		}

		// 2. 初始化角色
		roles, err := sd.seedRoles(tx)
		if err != nil {
			return err
		}

		// 3. 初始化角色权限关联
		if err := sd.seedRolePermissions(roles, permissions, tx); err != nil {
			return err
		}

		// 4. 初始化管理员用户
		if err := sd.seedAdminUser(roles, tx); err != nil {
			return err
		}

		// 5. 初始化菜单
		if err := sd.seedMenus(permissions, tx); err != nil {
			return err
		}

		log.Info("基础数据初始化完成")
		return nil
	})
}

func (sd *Seed) Check() (bool, error) {
	var count int64
	if err := sd.db.Model(&biz.User{}).Count(&count).Error; err != nil {
		// 访问数据库出错也跳过
		return true, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// seedPermissions 初始化权限数据
func (sd *Seed) seedPermissions(tx *gorm.DB) ([]*biz.Permission, error) {
	permissions := []*biz.Permission{
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "system",
			Alias:    "系统管理",
			Describe: "系统管理权限",
			Status:   1,
			Tags:     []string{"root", "admin"},
			Actions: []*biz.Action{
				{Key: "READ", Describe: "查看系统", Checked: true},
				{Key: "UPDATE", Describe: "更新系统", Checked: true},
			},
		},
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "user",
			Alias:    "用户管理",
			Describe: "用户管理权限",
			Status:   int64(administration.UserStatus_NORMAL),
			Tags:     []string{"root", "admin"},
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建用户", Checked: true},
				{Key: "READ", Describe: "查看用户", Checked: true},
				{Key: "UPDATE", Describe: "更新用户", Checked: true},
				{Key: "DELETE", Describe: "删除用户", Checked: true},
			},
		},
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "role",
			Alias:    "角色管理",
			Describe: "角色管理权限",
			Status:   1,
			Tags:     []string{"root", "admin"},
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建角色", Checked: true},
				{Key: "READ", Describe: "查看角色", Checked: true},
				{Key: "UPDATE", Describe: "更新角色", Checked: true},
				{Key: "DELETE", Describe: "删除角色", Checked: true},
			},
		},
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "permission",
			Alias:    "权限管理",
			Describe: "权限管理权限",
			Status:   1,
			Tags:     []string{"root", "admin"},
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建权限", Checked: true},
				{Key: "READ", Describe: "查看权限", Checked: true},
				{Key: "UPDATE", Describe: "更新权限", Checked: true},
				{Key: "DELETE", Describe: "删除权限", Checked: true},
			},
		},
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "menu",
			Alias:    "菜单管理",
			Describe: "菜单管理权限",
			Status:   1,
			Tags:     []string{"root", "admin"},
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建菜单", Checked: true},
				{Key: "READ", Describe: "查看菜单", Checked: true},
				{Key: "UPDATE", Describe: "更新菜单", Checked: true},
				{Key: "DELETE", Describe: "删除菜单", Checked: true},
			},
		},
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "crontab",
			Alias:    "定时任务",
			Describe: "定时任务管理权限",
			Status:   1,
			Tags:     []string{"root", "admin"},
			Actions: []*biz.Action{
				{Key: "CREATE", Describe: "创建定时任务", Checked: true},
				{Key: "READ", Describe: "查看定时任务", Checked: true},
				{Key: "UPDATE", Describe: "更新定时任务", Checked: true},
				{Key: "DELETE", Describe: "删除定时任务", Checked: true},
			},
		},
	}

	if err := tx.CreateInBatches(permissions, 100).Error; err != nil {
		log.Errorf("权限数据初始化失败: %v", err)
		return nil, err
	}

	log.Info("权限数据初始化完成")
	return permissions, nil
}

// seedRoles 初始化角色数据
func (sd *Seed) seedRoles(tx *gorm.DB) ([]*biz.Role, error) {
	roles := []*biz.Role{
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "root",
			Alias:    "超级管理员",
			Describe: "系统超级管理员，拥有所有权限",
			Status:   1,
		},
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "admin",
			Alias:    "管理员",
			Describe: "系统管理员，拥有几乎所有权限",
			Status:   1,
		},
		{
			UID:      sd.idGen.Generate().Int64(),
			Name:     "user",
			Alias:    "普通用户",
			Describe: "普通用户，拥有基本权限",
			Status:   1,
		},
	}

	if err := tx.CreateInBatches(roles, 100).Error; err != nil {
		log.Errorf("角色数据初始化失败: %v", err)
		return nil, err
	}
	log.Info("角色数据初始化完成")
	return roles, nil
}

// seedRolePermissions 初始化角色权限关联
func (sd *Seed) seedRolePermissions(roles []*biz.Role, permissions []*biz.Permission, tx *gorm.DB) error {
	bindings := make([]*biz.RolePermission, 0)
	for _, role := range roles {
		for _, permission := range permissions {
			if !slices.Contains(permission.Tags, role.Name) {
				continue
			}

			bindings = append(bindings, &biz.RolePermission{
				RoleID:     role.UID,
				PermID:     permission.UID,
				Actions:    permission.Actions,
				Name:       permission.Name,
				Alias:      permission.Alias,
				DataAccess: nil,
			})
		}
	}

	if err := tx.CreateInBatches(bindings, 100).Error; err != nil {
		log.Errorf("角色权限关联初始化失败: %v", err)
		return err
	}

	log.Info("角色权限关联初始化完成")
	return nil
}

// seedAdminUser 初始化管理员用户
func (sd *Seed) seedAdminUser(roles []*biz.Role, tx *gorm.DB) error {
	// 创建超级管理员用户
	adminUser := &biz.User{
		UID:      sd.idGen.Generate().Int64(),
		Username: "admin",
		Password: "12346578",
		Email:    "admin@example.com",
		Nickname: "超级管理员",
		Bio:      "系统超级管理员",
		Status:   int64(administration.UserStatus_NORMAL),
	}

	if err := tx.Create(adminUser).Error; err != nil {
		log.Errorf("创建管理员用户失败: %v", err)
		return err
	}

	rootRole, exist := lo.Find(roles, func(item *biz.Role) bool {
		return item.Name == "root"
	})
	if !exist {

	}

	// 为管理员用户分配超级管理员角色
	userRole := &biz.UserRole{
		UserID:    adminUser.UID,
		RoleID:    rootRole.UID, // 超级管理员角色
		CreatedAt: time.Now(),
	}
	if err := tx.Create(userRole).Error; err != nil {
		log.Errorf("管理员用户初始化失败: %v", err)
		return err
	}

	log.Info("管理员用户初始化完成")
	return nil
}

// seedMenus 初始化菜单数据
func (sd *Seed) seedMenus(permissions []*biz.Permission, tx *gorm.DB) error {
	menus := make([]*biz.Menu, 0)

	type MenuResource struct {
		Name string
		Icon string
		Path string
	}
	// resource
	mapResource := map[string]MenuResource{
		"dashboard": {
			Name: "仪表盘",
			Icon: "icon-app-box",
			Path: "/",
		},
		"system": {
			Name: "系统管理",
			Icon: "icon-config",
			Path: "/admin",
		},
		"user": {
			Name: "用户",
			Icon: "icon-user",
			Path: "/admin/user",
		},
		"role": {
			Name: "角色",
			Icon: "icon-addteam",
			Path: "/admin/role",
		},
		"permission": {
			Name: "权限",
			Icon: "icon-securityscan",
			Path: "/admin/permission",
		},
		"menu": {
			Name: "菜单",
			Icon: "icon-resource",
			Path: "/admin/menu",
		},
		"crontab": {
			Name: "定时任务",
			Icon: "icon-schedule",
			Path: "/admin/crontab",
		},
	}
	// dashboard
	dashboard := &biz.Menu{
		UID:          sd.idGen.Generate().Int64(),
		PID:          0,
		PermissionID: 0,
		Name:         "仪表盘",
		Icon:         "icon-app-box",
		Path:         "/",
		SortBy:       1,
		Hidden:       false,
		Status:       1,
	}
	menus = append(menus, dashboard)

	// system
	system := &biz.Menu{
		UID:          sd.idGen.Generate().Int64(),
		PID:          0,
		PermissionID: 0,
		Name:         "系统管理",
		Icon:         "icon-config",
		Path:         "/admin",
		SortBy:       2,
		Hidden:       false,
		Status:       1,
	}
	menus = append(menus, system)
	for i, permission := range permissions {
		if permission.Name == "system" {
			continue
		}

		if item, ok := mapResource[permission.Name]; ok {
			menus = append(menus, &biz.Menu{
				UID:          sd.idGen.Generate().Int64(),
				PID:          system.PID,
				PermissionID: permission.UID,
				Name:         item.Name,
				Icon:         item.Icon,
				Path:         item.Path,
				Hidden:       false,
				Status:       1,
				SortBy:       int64(i + 2),
			})
		}
	}
	if err := tx.CreateInBatches(menus, 100).Error; err != nil {
		log.Errorf("菜单数据初始化失败: %v", err)
		return err
	}
	log.Info("菜单数据初始化完成")
	return nil
}
