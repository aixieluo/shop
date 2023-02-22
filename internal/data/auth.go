package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "shop/api/user/v1"
	"shop/internal/biz"
)

type AuthRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	helper := log.NewHelper(log.With(logger, "module", "data/shop"))
	return &AuthRepo{
		data: data,
		log:  helper,
	}
}

func (a *AuthRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	user, err := a.data.uc.CreateUser(ctx, &v1.CreateUserRequest{
		Mobile:   u.Mobile,
		Password: u.Password,
		Nickname: u.Nickname,
		Email:    "",
	})
	if err != nil {
		return nil, err
	}
	bu := &biz.User{
		ID:       user.Id,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
	}
	return bu, nil
}

func (a *AuthRepo) UserById(ctx context.Context, id uint) (*biz.User, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthRepo) UserByMobile(ctx context.Context, mobile string) (*biz.User, error) {
	//TODO implement me
	panic("implement me")
}
