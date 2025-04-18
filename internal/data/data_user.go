package data

import (
	"context"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"

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

func (r *userRepo) SelectUserByID(ctx context.Context, id int64) (*biz.User, error) {
	return r.selectByField(ctx, "id", id)
}

func (r *userRepo) SelectList(ctx context.Context, pagination *protobuf.Pagination) ([]*biz.UserInfo, error) {
	var (
		list []*biz.UserInfo
		err  error
	)
	err = r.txm.WithContext(ctx).Model(&biz.User{}).
		Select("users.*", "group_concat(roles.id) as role_ids").
		Omit("users.password").
		Joins("LEFT JOIN users_bind_role ON users.id = users_bind_role.user_id").
		Joins("LEFT JOIN roles ON users_bind_role.role_id = roles.id").
		Group("users.id").
		Count(pagination.Count()).
		Scopes(pagination.Paginate()).
		Find(&list).Error
	return list, err
}

func (r *userRepo) BindNamespace(ctx context.Context, userID int64, namespaceID int64) error {
	return r.txm.WithContext(ctx).Create(&biz.UserNamespace{
		UserID:      userID,
		NamespaceID: namespaceID,
	}).Error
}

func (r *userRepo) BindRole(ctx context.Context, userID int64, roleID int) error {
	return r.txm.WithContext(ctx).Create(&biz.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *userRepo) UnbindNamespace(ctx context.Context, userID int64, namespaceID int64) error {
	return r.txm.WithContext(ctx).Where("user_id = ? AND namespace_id = ?", userID, namespaceID).Delete(&biz.UserNamespace{}).Error
}

func (r *userRepo) UnbindRole(ctx context.Context, userID int64, roleID int) error {
	return r.txm.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&biz.UserRole{}).Error
}

// Create implements biz.UserRepo.
func (r *userRepo) Create(ctx context.Context, user *biz.User) error {
	return r.txm.WithContext(ctx).Create(user).Error
}

// Delete implements biz.UserRepo.
func (r *userRepo) Delete(ctx context.Context, id int64) error {
	return r.txm.WithContext(ctx).Where("id = ?", id).Delete(&biz.User{}).Error
}

// Update implements biz.UserRepo.
func (r *userRepo) Update(ctx context.Context, id int64, user *biz.User) error {
	return r.txm.WithContext(ctx).Where("id = ?", id).Updates(user).Error
}
