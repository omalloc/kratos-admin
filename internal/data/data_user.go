package data

import (
	"context"
	"strconv"
	"strings"

	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"
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

func (r *userRepo) SelectUserByID(ctx context.Context, id int64) (*biz.UserInfo, error) {
	var ret biz.UserInfo

	err := r.txm.WithContext(ctx).Model(&biz.UserInfo{}).
		Select("users.*", "GROUP_CONCAT(roles.id) as role_ids").
		Omit("users.password").
		Joins("LEFT JOIN users_bind_role ON users.id = users_bind_role.user_id").
		Joins("LEFT JOIN roles ON users_bind_role.role_id = roles.id").
		Where("users.id = ?", id).
		First(&ret).Error

	return &ret, err
}

func (r *userRepo) SelectList(ctx context.Context, pagination *protobuf.Pagination) ([]*biz.UserInfo, error) {
	var (
		list []*biz.TempUserInfo
		err  error
	)
	err = r.txm.WithContext(ctx).Model(&biz.User{}).
		Select("users.*", "GROUP_CONCAT(roles.id) as role_ids").
		Omit("users.password").
		Joins("LEFT JOIN users_bind_role ON users.id = users_bind_role.user_id").
		Joins("LEFT JOIN roles ON users_bind_role.role_id = roles.id").
		Group("users.id").
		Count(pagination.Count()).
		Scopes(pagination.Paginate()).
		Find(&list).Error
	res := make([]*biz.UserInfo, 0, len(list))
	for i, item := range list {
		res = append(res, &biz.UserInfo{
			User: item.User,
		})
		if item.RoleIDs != "" {
			roleIDs := strings.Split(item.RoleIDs, ",")
			for _, roleID := range roleIDs {
				res[i].RoleIDs = append(res[i].RoleIDs, lo.Must(strconv.ParseInt(roleID, 10, 64)))
			}
		}
	}
	return res, err
}

func (r *userRepo) UpdateRole(ctx context.Context, userID int64, roleIDs []int64) error {
	err := r.txm.WithContext(ctx).Model(&biz.UserRole{}).Where("user_id = ?", userID).Delete(&biz.UserRole{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if len(roleIDs) > 0 {
		roles := make([]*biz.UserRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			roles = append(roles, &biz.UserRole{
				UserID: userID,
				RoleID: int(rid),
			})
		}
		if err := r.txm.WithContext(ctx).Create(&roles).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepo) BindRole(ctx context.Context, userID int64, roleID int) error {
	return r.txm.WithContext(ctx).Create(&biz.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
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
	return r.txm.WithContext(ctx).Model(&biz.User{}).
		Where("id = ?", id).
		Select("email", "nickname", "status").
		Updates(map[string]interface{}{
			"email":    user.Email,
			"nickname": user.Nickname,
			"status":   user.Status,
		}).Error
}
