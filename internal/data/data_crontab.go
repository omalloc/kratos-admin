package data

import (
	"context"
	"errors"
	"time"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
	"gorm.io/gorm"

	"github.com/omalloc/kratos-admin/internal/biz"
)

type crontabRepo struct {
	txm orm.Transaction
}

// NewCrontabRepo 创建定时任务仓储实现
func NewCrontabRepo(txm orm.Transaction) biz.CrontabRepo {
	return &crontabRepo{
		txm: txm,
	}
}

// Create 创建定时任务
func (r *crontabRepo) Create(ctx context.Context, crontab *biz.Crontab) error {
	return r.txm.WithContext(ctx).Create(crontab).Error
}

// Update 更新定时任务
func (r *crontabRepo) Update(ctx context.Context, uid int64, crontab *biz.Crontab) error {
	return r.txm.WithContext(ctx).Model(&biz.Crontab{}).
		Where("uid = ?", uid).
		Updates(crontab).Error
}

// Delete 删除定时任务
func (r *crontabRepo) Delete(ctx context.Context, uid int64) error {
	return r.txm.WithContext(ctx).Model(&biz.Crontab{}).Delete("uid = ?", uid).Error
}

// Get 获取定时任务详情
func (r *crontabRepo) Get(ctx context.Context, uid int64) (*biz.Crontab, error) {
	var crontab biz.Crontab
	err := r.txm.WithContext(ctx).Model(&biz.Crontab{}).
		Where("uid = ?", uid).First(&crontab).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回 nil 表示未找到，区别于其他错误
		}
		return nil, err
	}
	return &crontab, nil
}

// SelectList 获取定时任务列表（支持分页）
func (r *crontabRepo) SelectList(ctx context.Context, pagination *protobuf.Pagination) ([]*biz.Crontab, error) {
	var crontabs []*biz.Crontab

	tx := r.txm.WithContext(ctx).Model(&biz.Crontab{})

	// 执行分页查询
	if err := tx.
		Count(pagination.Count()).
		Scopes(pagination.Paginate()).
		Order("created_at DESC").
		Find(&crontabs).Error; err != nil {
		return nil, err
	}

	return crontabs, nil
}

// UpdateLastRunAt 更新最后执行时间
func (r *crontabRepo) UpdateLastRunAt(ctx context.Context, uid int64, lastrunAt time.Time) error {
	return r.txm.WithContext(ctx).Model(&biz.Crontab{}).
		Where("uid = ?", uid).
		Update("last_run_at", lastrunAt).Error
}

// GetByName 根据名称查询定时任务（用于重名检查）
func (r *crontabRepo) GetByName(ctx context.Context, name string) (*biz.Crontab, error) {
	var crontab biz.Crontab
	err := r.txm.WithContext(ctx).Where("name = ?", name).First(&crontab).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回 nil 表示未找到
		}
		return nil, err
	}
	return &crontab, nil
}
