package data

import (
	"context"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/kratos/orm/crud"

	"github.com/omalloc/kratos-admin/internal/biz"
)

type userRepo struct {
	crud.CRUD[biz.User]

	txm orm.Transaction
}

func NewUserRepo(txm orm.Transaction) biz.UserRepo {
	return &userRepo{
		CRUD: crud.New[biz.User](txm.WithContext(context.Background())),
		txm:  txm,
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

func (r *userRepo) BindNamespace(ctx context.Context, userID int64, namespaceID int64) error {
	return r.txm.WithContext(ctx).Create(&biz.UserNamespace{
		UserID:      userID,
		NamespaceID: namespaceID,
	}).Error
}

func (r *userRepo) BindRole(ctx context.Context, userID int64, roleID int64) error {
	return r.txm.WithContext(ctx).Create(&biz.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *userRepo) UnbindNamespace(ctx context.Context, userID int64, namespaceID int64) error {
	return r.txm.WithContext(ctx).Where("user_id = ? AND namespace_id = ?", userID, namespaceID).Delete(&biz.UserNamespace{}).Error
}

func (r *userRepo) UnbindRole(ctx context.Context, userID int64, roleID int64) error {
	return r.txm.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&biz.UserRole{}).Error
}
