package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "shop/app/user/api/user/v1"
	"shop/app/user/internal/biz"
)

type UserUsecase interface {
	List(ctx context.Context, page int, perPage int) ([]*biz.User, int, error)
	Get(ctx context.Context, id uint) (*biz.User, error)
	Create(ctx context.Context, user *biz.User) (*biz.User, error)
	Update(ctx context.Context, id uint, user *biz.User) (bool, error)
}

// UserService is a greeter service.
type UserService struct {
	v1.UnimplementedUserServer

	uc  UserUsecase
	log *log.Helper
}

// NewUserService new a greeter service.
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(logger)}
}

func (us UserService) ListUser(ctx context.Context, request *v1.ListUserRequest) (*v1.ListUserReply, error) {
	list, total, err := us.uc.List(ctx, int(request.Page), int(request.PrePage))
	if err != nil {
		return nil, err
	}
	rsp := &v1.ListUserReply{
		Total: int32(total),
	}
	for _, u := range list {
		rsp.Data = append(rsp.Data, &v1.GetUserReply{
			Id:       u.ID,
			Nickname: u.Nickname,
			Email:    u.Email,
			Mobile:   u.Mobile,
		})
	}
	return rsp, nil
}

func (us UserService) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserReply, error) {
	u, err := us.uc.Get(ctx, uint(request.Id))
	if err != nil {
		return nil, err
	}
	rsp := &v1.GetUserReply{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
		Mobile:   u.Mobile,
	}

	return rsp, nil
}

func (us UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	user, err := us.uc.Create(ctx, &biz.User{
		Mobile:   request.Mobile,
		Password: request.Password,
		Nickname: request.Nickname,
		Username: request.Username,
		Email:    request.Email,
	})
	if err != nil {
		return nil, err
	}

	userInfo := &v1.CreateUserReply{
		Id:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
	}

	return userInfo, nil
}

func (us UserService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*emptypb.Empty, error) {
	u, err := us.uc.Update(ctx, uint(request.Id), &biz.User{
		Mobile:   request.Mobile,
		Nickname: request.Nickname,
		Email:    request.Email,
	})
	if err != nil {
		return nil, err
	}
	if !u {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
