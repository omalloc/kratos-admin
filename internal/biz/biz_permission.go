package biz

import (
	"context"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
)

type Permission struct {
	ID       int64     `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"column:name;type:varchar(64);comment:权限名"`
	Alias    string    `json:"alias" gorm:"column:alias;type:varchar(64);comment:别名,展示名"`
	Describe string    `json:"describe" gorm:"column:describe;type:varchar(255);comment:描述"`
	Actions  []*Action `json:"actions" gorm:"column:actions;type:json;serializer:json;comment:操作"`
	Status   int64     `json:"status" gorm:"column:status;type:int;comment:状态"`

	orm.DBModel
}

type Action struct {
	Key      string `json:"key"`      // 动作名称,可自定义; CREATE, UPDATE, DELETE, READ
	Describe string `json:"describe"` // 描述
	Checked  bool   `json:"checked"`  // 默认选中
}

func (Permission) TableName() string {
	return "permissions"
}

type PermissionRepo interface {
	Create(ctx context.Context, permission *Permission) error
	Update(ctx context.Context, id int64, permission *Permission) error
	Delete(ctx context.Context, id int64) error
	GetPermission(ctx context.Context, id int64) (*Permission, error)
	SelectList(ctx context.Context, pagination *protobuf.Pagination) ([]*Permission, error)
}

type PermissionUsecase struct {
	txm orm.Transaction

	permissionRepo PermissionRepo
}

func NewPermissionUsecase(txm orm.Transaction, permissionRepo PermissionRepo) *PermissionUsecase {
	return &PermissionUsecase{txm: txm, permissionRepo: permissionRepo}
}

func (uc *PermissionUsecase) CreatePermission(ctx context.Context, permission *Permission) error {
	return uc.permissionRepo.Create(ctx, permission)
}

func (uc *PermissionUsecase) ListPermission(ctx context.Context, pagination *protobuf.Pagination) ([]*Permission, error) {
	return uc.permissionRepo.SelectList(ctx, pagination)
}

func (uc *PermissionUsecase) UpdatePermission(ctx context.Context, permission *Permission) error {
	return uc.permissionRepo.Update(ctx, permission.ID, permission)
}

func (uc *PermissionUsecase) DeletePermission(ctx context.Context, id int64) error {
	return uc.permissionRepo.Delete(ctx, id)
}
