package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
)

// Crontab 定时任务模型
type Crontab struct {
	ID        int64      `json:"id" gorm:"primaryKey;type:BIGINT;autoIncrement"`
	UID       int64      `json:"uid" gorm:"column:uid;type:BIGINT;uniqueIndex:idx_uid_uk"`
	Name      string     `json:"name" gorm:"column:name;type:varchar(255);comment:任务名称;not null"`
	Expr      string     `json:"expr" gorm:"column:expr;type:varchar(255);comment:Cron表达式;not null"`
	Action    string     `json:"action" gorm:"column:action;type:text;comment:任务动作;not null"`
	Describe  string     `json:"describe" gorm:"column:describe;type:varchar(500);comment:任务描述"`
	LastRunAt *time.Time `json:"last_run_at" gorm:"column:last_run_at;type:datetime;comment:上次执行时间"`

	orm.DBModel
}

// TableName 指定表名
func (Crontab) TableName() string {
	return "crontabs"
}

// CrontabRepo 定时任务仓储接口
type CrontabRepo interface {
	// 基础 CRUD 操作
	Create(ctx context.Context, crontab *Crontab) error
	Update(ctx context.Context, id int64, crontab *Crontab) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*Crontab, error)

	// 列表查询
	SelectList(ctx context.Context, pagination *protobuf.Pagination) ([]*Crontab, error)

	// 更新最后执行时间
	UpdateLastRunAt(ctx context.Context, id int64, lastrunAt time.Time) error

	// 按名称查询（用于重名检查）
	GetByName(ctx context.Context, name string) (*Crontab, error)
}

// CrontabUsecase 定时任务用例
type CrontabUsecase struct {
	log         *log.Helper
	txm         orm.Transaction
	crontabRepo CrontabRepo
}

// NewCrontabUsecase 创建定时任务用例
func NewCrontabUsecase(repo CrontabRepo, txm orm.Transaction, logger log.Logger) *CrontabUsecase {
	return &CrontabUsecase{
		log:         log.NewHelper(logger),
		txm:         txm,
		crontabRepo: repo,
	}
}

// CreateCrontab 创建定时任务
func (uc *CrontabUsecase) CreateCrontab(ctx context.Context, crontab *Crontab) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		// 检查名称是否重复
		existing, err := uc.crontabRepo.GetByName(ctx, crontab.Name)
		if err == nil && existing != nil {
			return errors.New(400, "CRONTAB_NAME_EXISTS", "定时任务名称已存在")
		}

		return uc.crontabRepo.Create(ctx, crontab)
	})
}

// UpdateCrontab 更新定时任务
func (uc *CrontabUsecase) UpdateCrontab(ctx context.Context, crontab *Crontab) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		// 检查是否存在
		existing, err := uc.crontabRepo.Get(ctx, crontab.ID)
		if err != nil {
			return err
		}
		if existing == nil {
			return errors.New(404, "CRONTAB_NOT_FOUND", "定时任务不存在")
		}

		// 如果修改了名称，检查新名称是否重复
		if existing.Name != crontab.Name {
			nameExists, err := uc.crontabRepo.GetByName(ctx, crontab.Name)
			if err == nil && nameExists != nil {
				return errors.New(400, "CRONTAB_NAME_EXISTS", "定时任务名称已存在")
			}
		}

		return uc.crontabRepo.Update(ctx, crontab.ID, crontab)
	})
}

// DeleteCrontab 删除定时任务
func (uc *CrontabUsecase) DeleteCrontab(ctx context.Context, id int64) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		// 检查是否存在
		existing, err := uc.crontabRepo.Get(ctx, id)
		if err != nil {
			return err
		}
		if existing == nil {
			return errors.New(404, "CRONTAB_NOT_FOUND", "定时任务不存在")
		}

		return uc.crontabRepo.Delete(ctx, id)
	})
}

// GetCrontab 获取定时任务详情
func (uc *CrontabUsecase) GetCrontab(ctx context.Context, id int64) (*Crontab, error) {
	crontab, err := uc.crontabRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if crontab == nil {
		return nil, errors.New(404, "CRONTAB_NOT_FOUND", "定时任务不存在")
	}

	return crontab, nil
}

// ListCrontab 获取定时任务列表
func (uc *CrontabUsecase) ListCrontab(ctx context.Context, pagination *protobuf.Pagination) ([]*Crontab, error) {
	return uc.crontabRepo.SelectList(ctx, pagination)
}

// UpdateLastrunAt 更新最后执行时间
func (uc *CrontabUsecase) UpdateLastrunAt(ctx context.Context, id int64, lastrunAt time.Time) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		return uc.crontabRepo.UpdateLastRunAt(ctx, id, lastrunAt)
	})
}
