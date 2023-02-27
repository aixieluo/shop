package biz

import (
	"context"

	v1 "shop/app/user/api/user/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// User is a User model.
type User struct {
	ID       uint64 `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Email    string `json:"email,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Password string `json:"password,omitempty"`
}

// UserRepo is a Greater repo.
type UserRepo interface {
	List(ctx context.Context, page int, perPage int) ([]*User, int, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	Create(ctx context.Context, bu *User) (*User, error)
	Update(ctx context.Context, id uint, user *User) (bool, error)
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc UserUsecase) List(ctx context.Context, page int, perPage int) ([]*User, int, error) {
	return uc.repo.List(ctx, page, perPage)
}

func (uc UserUsecase) Get(ctx context.Context, id uint) (*User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc UserUsecase) Create(ctx context.Context, user *User) (*User, error) {
	//uc.log.WithContext(ctx).Infof("Create: %v", uc)
	return uc.repo.Create(ctx, user)
}

func (uc UserUsecase) Update(ctx context.Context, id uint, user *User) (bool, error) {
	return uc.repo.Update(ctx, id, user)
}
