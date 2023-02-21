package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "shop/api/shop/v1"
	"shop/internal/biz"
	"shop/internal/pkg/captcha"
)

type AuthUsecase interface {
	Register(ctx context.Context, ba *biz.User, bc *biz.Captcha) (string, error)
	Login(ctx context.Context, ba biz.User) (string, error)
	GetCaptcha(ctx context.Context) (*captcha.Info, error)
}

type AuthService struct {
	v1.UnimplementedAuthServer
	uc  AuthUsecase
	log *log.Helper
}

func NewAuthService(uc *biz.AuthUsecase, logger log.Logger) *AuthService {
	return &AuthService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/shop")),
	}
}

func (s *AuthService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	token, err := s.uc.Register(ctx, &biz.User{
		Password: req.Password,
		Mobile:   req.Mobile,
	}, &biz.Captcha{
		CaptchaId: req.CaptchaId,
		Captcha:   req.Captcha,
	})
	if err != nil {
		return nil, err
	}
	return &v1.RegisterReply{Token: token}, nil
}
func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	token, err := s.uc.Login(ctx, biz.User{Mobile: req.Mobile, Password: req.Password})
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{Token: token}, nil
}
func (s *AuthService) Captcha(ctx context.Context, req *emptypb.Empty) (*v1.CaptchaReply, error) {
	info, err := s.uc.GetCaptcha(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CaptchaReply{
		CaptchaId:  info.CaptchaId,
		CaptchaPic: info.PicPath,
	}, nil
}
func (s *AuthService) ModifyPassword(ctx context.Context, req *v1.ModifyPasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
