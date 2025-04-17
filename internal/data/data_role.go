package data

import (
	"context"
	"errors"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/kratos/orm/crud"
	"github.com/omalloc/contrib/protobuf"
	"gorm.io/gorm"

	"github.com/omalloc/kratos-admin/internal/biz"
)

type roleRepo struct {
	crud.CRUD[biz.Role]

	txm orm.Transaction
}

func NewRoleRepo(txm orm.Transaction) biz.RoleRepo {
	return &roleRepo{
		CRUD: crud.New[biz.Role](txm.WithContext(context.Background())),
		txm:  txm,
	}
}

func (r *roleRepo) Exist(ctx context.Context, name string) bool {
	err := r.txm.WithContext(ctx).Model(&biz.Role{}).Where("name = ?", name).First(&biz.Role{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (r *roleRepo) SelectByName(ctx context.Context, name string) (*biz.Role, error) {
	var role biz.Role
	err := r.txm.WithContext(ctx).Model(&biz.Role{}).Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *roleRepo) SelectFilterList(ctx context.Context, pagination *protobuf.Pagination) ([]*biz.Role, error) {
	var (
		list []*biz.Role
		err  error
	)
	err = r.txm.WithContext(ctx).Model(&biz.Role{}).
		Count(pagination.Count()).
		Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Find(&list).Error

	return list, err
}

// BindPermission implements biz.RoleRepo.
func (r *roleRepo) BindPermission(ctx context.Context, roleID int64, permissionID int64, actions []*biz.Action, dataAccess []*biz.Action) error {
	return r.txm.WithContext(ctx).Model(&biz.RolePermission{}).
		Create(&biz.RolePermission{
			RoleID:     roleID,
			PermID:     permissionID,
			Actions:    actions,
			DataAccess: dataAccess,
		}).Error
}

// UnbindPermission implements biz.RoleRepo.
func (r *roleRepo) UnbindPermission(ctx context.Context, roleID int64, permissionID int64) error {
	return r.txm.WithContext(ctx).Model(&biz.RolePermission{}).
		Where("role_id = ? AND perm_id = ?", roleID, permissionID).
		Delete(&biz.RolePermission{}).Error
}
