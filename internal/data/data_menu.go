package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
	"github.com/omalloc/kratos-admin/internal/biz"
)

type menuRepo struct {
	txm orm.Transaction
	log *log.Helper
}

// NewMenuRepo .
func NewMenuRepo(txm orm.Transaction, logger log.Logger) biz.MenuRepo {
	return &menuRepo{
		txm: txm,
		log: log.NewHelper(logger),
	}
}

// Create 创建菜单
func (r *menuRepo) Create(ctx context.Context, m *biz.Menu) error {
	return r.txm.WithContext(ctx).Create(m).Error
}

// Update 更新菜单
func (r *menuRepo) Update(ctx context.Context, m *biz.Menu) error {
	return r.txm.WithContext(ctx).Model(&biz.Menu{}).
		Where("uid = ?", m.UID).
		Updates(m).Error
}

// Delete 删除菜单
func (r *menuRepo) Delete(ctx context.Context, uid int64) error {
	return r.txm.WithContext(ctx).Model(&biz.Menu{}).Delete("uid = ?", uid).Error
}

// Get 获取菜单
func (r *menuRepo) Get(ctx context.Context, uid int64) (*biz.Menu, error) {
	var m biz.Menu
	if err := r.txm.WithContext(ctx).Model(&biz.Menu{}).
		Where("uid = ?", uid).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// SelectList 获取菜单列表
func (r *menuRepo) SelectList(ctx context.Context, pagination *protobuf.Pagination, name string, status int32) ([]*biz.Menu, error) {
	var menus []*biz.Menu
	query := r.txm.WithContext(ctx).Model(&biz.Menu{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("sort_by ASC").
		Count(pagination.Count()).
		Scopes(pagination.Paginate()).
		Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepo) SelectAll(ctx context.Context) ([]*biz.Menu, error) {
	var menus []*biz.Menu
	if err := r.txm.WithContext(ctx).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
