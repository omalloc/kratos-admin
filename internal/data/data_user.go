package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
	"gorm.io/gorm"

	"github.com/omalloc/kratos-admin/internal/biz"
)

type userRepo struct {
	txm orm.Transaction
}

func NewUserRepo(txm orm.Transaction) biz.UserRepo {
	return &userRepo{
		txm: txm,
	}
}

func (r *userRepo) selectByField(ctx context.Context, field string, val any) (*biz.User, error) {
	var (
		user biz.User
		err  error
	)

	err = r.txm.WithContext(ctx).
		Where(field, val).First(&user).Error
	return &user, err
}

func (r *userRepo) SelectUserByName(ctx context.Context, name string) (*biz.User, error) {
	return r.selectByField(ctx, "username", name)
}

func (r *userRepo) SelectUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	return r.selectByField(ctx, "email", email)
}

// SelectUserByNameOrEmail implements biz.UserRepo.
func (r *userRepo) SelectUserByNameOrEmail(ctx context.Context, value string) (*biz.User, error) {
	var (
		user biz.User
		err  error
	)

	err = r.txm.WithContext(ctx).
		Where(&biz.User{Username: value}).
		Or(&biz.User{Email: value}).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *userRepo) SelectUserByUID(ctx context.Context, uid int64) (*biz.UserInfo, error) {
	var ret biz.UserInfo

	tx := r.txm.WithContext(ctx)

	subQuery := tx.Model(&biz.UserRole{}).
		Select("user_id, GROUP_CONCAT(role_id) AS role_ids").
		Where("user_id = ?", uid).
		Group("user_id")

	err := tx.Model(&biz.UserInfo{}).
		Select(
			"users.uid",
			"users.username",
			"users.email",
			"users.avatar_id",
			"users.nickname",
			"users.bio",
			"users.status",
			"users.created_at",
			"users.updated_at",
			"r.role_ids",
		).
		Joins("LEFT JOIN (?) AS r ON users.uid = r.user_id", subQuery).
		Where("users.uid = ?", uid).
		Order("users.uid DESC").
		Take(&ret).Error

	return &ret, err
}

func (r *userRepo) SelectList(ctx context.Context, pagination *protobuf.Pagination, filter *biz.UserQueryFilter) ([]*biz.UserInfo, error) {
	var (
		list []*biz.UserInfo
		err  error
	)
	tx := r.txm.WithContext(ctx).Model(&biz.User{}).
		Select("users.uid",
			"users.username",
			"users.email",
			"users.avatar_id",
			"users.nickname",
			"users.bio",
			"users.status",
			"users.last_login",
			"GROUP_CONCAT(roles.uid) as role_ids").
		Omit("users.password").
		Joins("LEFT JOIN users_bind_role ON users.uid = users_bind_role.user_id").
		Joins("LEFT JOIN roles ON users_bind_role.role_id = roles.uid")
	if filter != nil {
		if filter.Status > 0 {
			tx = tx.Where("users.status = ?", filter.Status)
		}
		if strings.Contains(filter.Username, "@") {
			tx = tx.Where("users.email LIKE ?", fmt.Sprintf("%%%s%%", filter.Username))
		} else if filter.Username != "" {
			tx = tx.Where("users.username LIKE ?", fmt.Sprintf("%%%s%%", filter.Username))
		}
	}

	err = tx.
		Group("users.uid, users.username, users.email, users.avatar_id, users.nickname, users.bio, users.status, users.last_login").
		Count(pagination.Count()).
		Scopes(pagination.Paginate()).
		Find(&list).Error
	return list, err
}

func (r *userRepo) UpdateRole(ctx context.Context, userID int64, roleIDs []int64) error {
	err := r.txm.WithContext(ctx).
		Model(&biz.UserRole{}).
		Where("user_id = ?", userID).
		Delete(&biz.UserRole{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if len(roleIDs) > 0 {
		roles := make([]*biz.UserRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			roles = append(roles, &biz.UserRole{
				UserID: userID,
				RoleID: rid,
			})
		}
		if err := r.txm.WithContext(ctx).Create(&roles).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepo) BindRole(ctx context.Context, userID int64, roleID int64) error {
	return r.txm.WithContext(ctx).Create(&biz.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *userRepo) UnbindRole(ctx context.Context, userID int64, roleID int64) error {
	return r.txm.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&biz.UserRole{}).Error
}

// Create implements biz.UserRepo.
func (r *userRepo) Create(ctx context.Context, user *biz.User) error {
	return r.txm.WithContext(ctx).Create(user).Error
}

// Delete implements biz.UserRepo.
func (r *userRepo) Delete(ctx context.Context, uid int64) error {
	return r.txm.WithContext(ctx).Where("uid = ?", uid).Delete(&biz.User{}).Error
}

// Update implements biz.UserRepo.
func (r *userRepo) Update(ctx context.Context, uid int64, user *biz.User) error {
	return r.txm.WithContext(ctx).Model(&biz.User{}).
		Where("uid = ?", uid).
		Updates(user).Error
}

func (r *userRepo) UpdateStatus(ctx context.Context, uid int64, status int64) error {
	return r.txm.WithContext(ctx).Model(&biz.User{}).
		Where("uid = ?", uid).
		Update("status", status).Error
}

func (r *userRepo) UpdateLastLogin(ctx context.Context, uid int64) error {
	return r.txm.WithContext(ctx).
		Model(&biz.User{}).
		Where("uid = ?", uid).
		Update("last_login", time.Now()).Error
}
