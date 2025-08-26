package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/protobuf"
	"gorm.io/gorm"
)

type Menu struct {
	ID           int64     `json:"id" gorm:"primaryKey;type:BIGINT;autoIncrement"`
	UID          int64     `json:"uid" gorm:"column:uid;type:BIGINT;uniqueIndex:idx_uid_uk"`
	PID          int64     `json:"pid" gorm:"column:pid;type:BIGINT;comment:父级ID"` // zero is root node.
	PermissionID int64     `json:"permission_id" gorm:"column:permission_id;type:BIGINT;comment:权限ID"`
	Name         string    `json:"name" gorm:"column:name;type:varchar(255);comment:名称"`
	Icon         string    `json:"icon" gorm:"column:icon;type:varchar(255);comment:图标"`
	Path         string    `json:"path" gorm:"column:path;type:varchar(255);comment:路径"`
	SortBy       int64     `json:"sort_by" gorm:"column:sort_by;type:int;comment:排序"`
	Hidden       bool      `json:"hidden" gorm:"column:hidden;type:tinyint;comment:是否隐藏=0显示1隐藏"`
	Status       int32     `json:"status" gorm:"column:status;type:int;comment:状态"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;type:datetime;comment:创建时间"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;comment:更新时间"`
	gorm.DeletedAt
}

func (Menu) TableName() string {
	return "menus"
}

type MenuRepo interface {
	Create(ctx context.Context, m *Menu) error
	Update(ctx context.Context, m *Menu) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*Menu, error)
	SelectList(ctx context.Context, pagination *protobuf.Pagination, name string, status int32) ([]*Menu, error)
	SelectAll(ctx context.Context) ([]*Menu, error)
}

type MenuUsecase struct {
	repo MenuRepo
	log  *log.Helper
}

func NewMenuUsecase(repo MenuRepo, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Create 创建菜单
func (uc *MenuUsecase) Create(ctx context.Context, m *Menu) error {
	if m.Name == "" {
		return errors.New(400, "MENU_NAME_EMPTY", "菜单名称不能为空")
	}

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	return uc.repo.Create(ctx, m)
}

// Update 更新菜单
func (uc *MenuUsecase) Update(ctx context.Context, m *Menu) error {
	if m.ID <= 0 {
		return errors.New(400, "MENU_ID_INVALID", "菜单ID无效")
	}

	if m.Name == "" {
		return errors.New(400, "MENU_NAME_EMPTY", "菜单名称不能为空")
	}

	m.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, m)
}

// Delete 删除菜单
func (uc *MenuUsecase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New(400, "MENU_ID_INVALID", "菜单ID无效")
	}

	return uc.repo.Delete(ctx, id)
}

// SelectByID 获取菜单
func (uc *MenuUsecase) SelectByID(ctx context.Context, id int64) (*Menu, error) {
	if id <= 0 {
		return nil, errors.New(400, "MENU_ID_INVALID", "菜单ID无效")
	}

	return uc.repo.Get(ctx, id)
}

// SelectList 获取菜单列表
func (uc *MenuUsecase) SelectList(ctx context.Context, pagination *protobuf.Pagination, name string, status int32) ([]*Menu, error) {
	return uc.repo.SelectList(ctx, pagination, name, status)
}

func (uc *MenuUsecase) SelectAll(ctx context.Context) ([]*Menu, error) {
	return uc.repo.SelectAll(ctx)
}
