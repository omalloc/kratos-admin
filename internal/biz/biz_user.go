package biz

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	passportpb "github.com/omalloc/kratos-admin/api/console/passport"
)

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"column:username;type:varchar(64);comment:用户名"`
	Password string `json:"-" gorm:"column:password;type:varchar(64);comment:密码"`
	Email    string `json:"email" gorm:"column:email;type:varchar(64);comment:邮箱"`
	Nickname string `json:"nickname" gorm:"column:nickname;type:varchar(64);comment:昵称"`
	Bio      string `json:"bio" gorm:"column:bio;type:varchar(255);comment:个人简介"`
	AvatarID int64  `json:"avatar_id" gorm:"column:avatar_id;comment:头像"`
	Status   int64  `json:"status" gorm:"column:status;type:int;comment:状态"` // 0: 禁用, 1: 正常

	orm.DBModel
}

type UserInfo struct {
	User

	RoleIDs []int64 `json:"role_ids" gorm:"serializer:intslice"`
}

type TempUserInfo struct {
	User

	RoleIDs string `json:"role_ids" `
}

type UserRoleInfo struct {
	User

	Roles []*Role `json:"roles"`
}

func (User) TableName() string {
	return "users"
}

type UserRole struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"column:user_id;type:int;uniqueIndex:idx_unique_user_role;comment:用户ID"`
	RoleID    int       `json:"role_id" gorm:"column:role_id;type:int;uniqueIndex:idx_unique_user_role;comment:角色ID"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime;comment:创建时间"` // 授权时间
}

func (UserRole) TableName() string {
	return "users_bind_role"
}

type UserRepo interface {
	SelectList(ctx context.Context, pagination *protobuf.Pagination) ([]*UserInfo, error)
	SelectUserByID(ctx context.Context, id int64) (*UserInfo, error)
	SelectUserByName(ctx context.Context, name string) (*User, error)
	SelectUserByEmail(ctx context.Context, email string) (*User, error)
	SelectUserByNameOrEmail(ctx context.Context, value string) (*User, error)

	BindRole(ctx context.Context, userID int64, roleID int) error
	UnbindRole(ctx context.Context, userID int64, roleID int) error
	UpdateRole(ctx context.Context, userID int64, roleIDs []int64) error

	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id int64, user *User) error
	Delete(ctx context.Context, id int64) error
}

type UserUsecase struct {
	log      *log.Helper
	txm      orm.Transaction
	userRepo UserRepo
	roleRepo RoleRepo
}

func NewUserUsecase(repo UserRepo, roleRepo RoleRepo, txm orm.Transaction, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: repo,
		roleRepo: roleRepo,
		txm:      txm,
		log:      log.NewHelper(logger),
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, user *User) error {
	hp, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hp)

	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		curr, err1 := uc.userRepo.SelectUserByName(ctx, user.Username)
		if err1 != nil && !errors.Is(err1, gorm.ErrRecordNotFound) {
			uc.log.Errorf("SelectUserByName error: %v", err1)
			return errors.New(400, "USER_EXIST", "用户存在")
		}

		if curr.ID > 0 {
			return errors.New(400, "USER_EXIST", "用户存在")
		}

		return uc.userRepo.Create(ctx, user)
	})
}

// GetUser 获取用户信息
func (uc *UserUsecase) GetUser(ctx context.Context, id int64) (*UserRoleInfo, error) {
	var (
		ret UserRoleInfo
		err error
	)

	err = uc.txm.Transaction(ctx, func(ctx context.Context) error {
		user, err := uc.userRepo.SelectUserByID(ctx, id)
		if err != nil {
			return err
		}

		roles, err := uc.roleRepo.SelectRolePermission(ctx, user.RoleIDs)
		if err != nil {
			return err
		}

		ret.User = user.User
		ret.Roles = lo.Map(roles, func(item *RoleJoinPermission, _ int) *Role {
			return &Role{
				ID:          item.ID,
				Name:        item.Name,
				Describe:    item.Describe,
				Status:      item.Status,
				Permissions: item.Permissions,
			}
		})

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// UpdateUser 更新用户信息
func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) error {
	if user.ID <= 0 {
		return errors.New(400, "INVALID_USER_ID", "无效的用户ID")
	}
	return uc.userRepo.Update(ctx, user.ID, user)
}

// DeleteUser 删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New(400, "INVALID_USER_ID", "无效的用户ID")
	}
	return uc.userRepo.Delete(ctx, id)
}

// ListUser 获取用户列表
func (uc *UserUsecase) ListUser(ctx context.Context, pagination *protobuf.Pagination) ([]*UserInfo, error) {
	return uc.userRepo.SelectList(ctx, pagination)
}

func (uc *UserUsecase) BindRole(ctx context.Context, userID int64, roleID int) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		return uc.userRepo.BindRole(ctx, userID, roleID)
	})
}

func (uc *UserUsecase) UnbindRole(ctx context.Context, userID int64, roleID int) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		return uc.userRepo.UnbindRole(ctx, userID, roleID)
	})
}

func (uc *UserUsecase) UpdateRole(ctx context.Context, userID int64, roleIDs []int64) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		return uc.userRepo.UpdateRole(ctx, userID, roleIDs)
	})
}

func (uc *UserUsecase) Login(ctx context.Context, username string, password string, autoLogin bool) (*User, error) {
	user, err := uc.userRepo.SelectUserByNameOrEmail(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, passportpb.ErrorUserOrPasswordError("")
		}
		return nil, err
	}

	if user.Password == "" {
		return nil, passportpb.ErrorUserOrPasswordError("")
	}

	// 密码前缀检查，如果非加密密码 进行一次加密并保存
	if !strings.HasPrefix(user.Password, "$2a$") {
		buf, err1 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err1 != nil {
			return nil, err1
		}
		user.Password = string(buf)

		// 更新一次加密密码
		if err1 := uc.userRepo.Update(ctx, user.ID, user); err1 != nil {
			return nil, err1
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New(400, "USER_OR_PASSWORD_ERROR", "用户名或密码错误")
	}

	if user.Status == 0 {
		return nil, errors.New(400, "USER_DISABLED", "用户已禁用")
	}

	return user, nil
}

func (uc *UserUsecase) UpdatePassword(ctx context.Context, email string, password string) error {
	return uc.txm.Transaction(ctx, func(ctx context.Context) error {
		user, err := uc.userRepo.SelectUserByEmail(ctx, email)
		if err != nil {
			return err
		}

		buf, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newPassword := string(buf)
		if user.Password == newPassword {
			return errors.New(400, "NEW_PASSWORD_SAME_AS_OLD", "新密码与旧密码相同")
		}

		user.Password = newPassword
		return uc.userRepo.Update(ctx, user.ID, user)
	})
}
