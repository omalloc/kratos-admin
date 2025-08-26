package data

import (
	"context"
	"errors"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/kratos/orm/crud"
	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

func (r *roleRepo) SelectID(ctx context.Context, uid int64) (*biz.RoleJoinPermission, error) {
	var ret *biz.RoleJoinPermission

	err := r.txm.WithContext(ctx).Model(&biz.RoleJoinPermission{}).
		Where("roles.uid = ?", uid).
		Preload(clause.Associations).
		Find(&ret).Error

	return ret, err
}

// func (r *roleRepo) SelectFilterList(ctx context.Context, pagination *protobuf.Pagination) ([]*biz.Role, error) {
// 	var (
// 		list []*biz.Role
// 		err  error
// 	)
// 	err = r.txm.WithContext(ctx).Model(&biz.Role{}).
// 		Count(pagination.Count()).
// 		Offset(pagination.Offset()).
// 		Limit(pagination.Limit()).
// 		Find(&list).Error

// 	return list, err
// }

func (r *roleRepo) SelectFilterList(ctx context.Context, pagination *protobuf.Pagination) ([]*biz.RoleJoinPermission, error) {
	var (
		ret []*biz.RoleJoinPermission
		err error
	)

	err = r.txm.WithContext(ctx).Model(&biz.RoleJoinPermission{}).
		Count(pagination.Count()).
		Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Preload(clause.Associations).
		Find(&ret).Error

	return ret, err
}

func (r *roleRepo) SelectByUserID(ctx context.Context, userID int64) ([]*biz.Role, error) {
	var roles []*biz.Role

	err := r.txm.WithContext(ctx).Model(&biz.Role{}).
		Joins("LEFT JOIN users_bind_role ON roles.uid = users_bind_role.role_id").
		Where("users_bind_role.user_id = ?", userID).
		Find(&roles).Error

	return roles, err
}

func (r *roleRepo) SelectRolePermission(ctx context.Context, roleIDs []int64) ([]*biz.RoleJoinPermission, error) {
	var (
		roles           []*biz.Role
		rolePermissions []*biz.RolePermission
	)

	// query roles

	txm := r.txm.WithContext(ctx)
	err := txm.Model(&biz.Role{}).
		Where("roles.uid IN (?)", roleIDs).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// query role_permissions
	err = txm.Model(&biz.RolePermission{}).
		Select("roles_bind_permission.*, permissions.name, permissions.alias").
		Joins("LEFT JOIN permissions ON roles_bind_permission.perm_id = permissions.uid").
		Where("roles_bind_permission.role_id IN (?)", roleIDs).
		Find(&rolePermissions).Error

	ret := lo.Map(roles, func(role *biz.Role, _ int) *biz.RoleJoinPermission {
		return &biz.RoleJoinPermission{
			Role: *role,
			Permissions: lo.Filter(rolePermissions, func(item *biz.RolePermission, _ int) bool {
				return item.RoleID == role.UID
			}),
		}
	})
	return ret, err
}

// BindPermission implements biz.RoleRepo.
func (r *roleRepo) BindPermission(ctx context.Context, roleID int64, permissionID int64, actions []*biz.Action, dataAccess []*biz.Action) error {
	var rolePerm biz.RolePermission
	result := r.txm.WithContext(ctx).Model(&biz.RolePermission{}).
		Where("role_id = ? AND perm_id = ?", roleID, permissionID).
		First(&rolePerm)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new record if not found
			return r.txm.WithContext(ctx).Model(&biz.RolePermission{}).
				Create(&biz.RolePermission{
					RoleID:     roleID,
					PermID:     permissionID,
					Actions:    actions,
					DataAccess: dataAccess,
				}).Error
		}
		return result.Error
	}

	// Update existing record
	rolePerm.Actions = actions
	rolePerm.DataAccess = dataAccess
	return r.txm.WithContext(ctx).Save(&rolePerm).Error
}

// UnbindPermission implements biz.RoleRepo.
func (r *roleRepo) UnbindPermission(ctx context.Context, roleID int64, permissionID int64) error {
	return r.txm.WithContext(ctx).Model(&biz.RolePermission{}).
		Where("role_id = ? AND perm_id = ?", roleID, permissionID).
		Delete(&biz.RolePermission{}).Error
}

func (r *roleRepo) Update(ctx context.Context, uid int64, role *biz.Role) error {
	return r.txm.WithContext(ctx).Where("uid = ?", uid).Updates(role).Error
}

func (r *roleRepo) GetAll(ctx context.Context) ([]*biz.Role, error) {
	var roles []*biz.Role
	return roles, r.txm.WithContext(ctx).Model(&biz.Role{}).Find(&roles).Error
}
