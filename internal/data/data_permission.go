package data

import (
	"context"
	"fmt"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"

	"github.com/omalloc/kratos-admin/internal/biz"
)

type permissionRepo struct {
	txm orm.Transaction
}

func NewPermissionRepo(txm orm.Transaction) biz.PermissionRepo {
	return &permissionRepo{txm: txm}
}

func (r *permissionRepo) Create(ctx context.Context, permission *biz.Permission) error {
	return r.txm.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepo) Update(ctx context.Context, id int64, permission *biz.Permission) error {
	return r.txm.WithContext(ctx).Where("id = ?", id).Updates(permission).Error
}

func (r *permissionRepo) Delete(ctx context.Context, id int64) error {
	return r.txm.WithContext(ctx).Where("id = ?", id).Delete(&biz.Permission{}).Error
}

func (r *permissionRepo) GetPermission(ctx context.Context, id int64) (*biz.Permission, error) {
	var permission biz.Permission
	if err := r.txm.WithContext(ctx).Where("id = ?", id).First(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *permissionRepo) SelectList(ctx context.Context, name string, status int32, pagination *protobuf.Pagination) ([]*biz.Permission, error) {
	var permissions []*biz.Permission

	tx := r.txm.WithContext(ctx).Model(&biz.Permission{})
	if name != "" {
		tx = tx.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if status != 0 {
		tx = tx.Where("status = ?", status)
	}

	if err := tx.
		Count(pagination.Count()).
		Scopes(pagination.Paginate()).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *permissionRepo) SelectAll(ctx context.Context, scoped bool) ([]*biz.Permission, error) {
	var permissions []*biz.Permission

	tx := r.txm.WithContext(ctx).Model(&biz.Permission{})
	if scoped {
		// tx = tx.Where("status = ?", 1)
	}

	if err := tx.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}
