package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
	"gorm.io/gorm"

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
	return r.txm.WithContext(ctx).Model(&biz.Menu{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"pid":        m.PID,
		"name":       m.Name,
		"icon":       m.Icon,
		"path":       m.Path,
		"sort_by":    m.SortBy,
		"hidden":     m.Hidden,
		"status":     m.Status,
		"updated_at": m.UpdatedAt,
	}).Error
}

// Delete 删除菜单
func (r *menuRepo) Delete(ctx context.Context, id int64) error {
	return r.txm.WithContext(ctx).Delete(&biz.Menu{}, id).Error
}

// Get 获取菜单
func (r *menuRepo) Get(ctx context.Context, id int64) (*biz.Menu, error) {
	var m biz.Menu
	if err := r.txm.WithContext(ctx).First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
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
