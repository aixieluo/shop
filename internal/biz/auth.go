package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"time"

	"shop/internal/conf"
	"shop/internal/pkg/captcha"
	"shop/internal/pkg/middleware/auth"
)

type User struct {
	ID        uint64
	Nickname  string
	Mobile    string
	Password  string
	CreatedAt time.Time
}

type Captcha struct {
	CaptchaId string
	Captcha   string
}

type AuthRepo interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	UserById(ctx context.Context, id uint) (*User, error)
	UserByMobile(ctx context.Context, mobile string) (*User, error)
}

type AuthUsecase struct {
	uRepo      AuthRepo
	log        *log.Helper
	signingKey string
}

func NewAuthUsecase(uRepo AuthRepo, logger log.Logger, c *conf.Auth) *AuthUsecase {
	return &AuthUsecase{
		uRepo:      uRepo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/shop")),
		signingKey: c.JwtKey,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, ba *User, bc *Captcha) (string, error) {
	user, err := uc.uRepo.CreateUser(ctx, ba)
	if err != nil {
		return "", err
	}
	claims := auth.CustomClaims{
		ID:       user.ID,
		Nickname: user.Nickname,
		RegisteredClaims: jwtV4.RegisteredClaims{
			NotBefore: jwtV4.NewNumericDate(time.Now()),
			ExpiresAt: jwtV4.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			Issuer:    "Gyl",
		},
	}
	token, err := auth.CreateToken(claims, uc.signingKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, ba User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (uc *AuthUsecase) GetCaptcha(ctx context.Context) (*captcha.Info, error) {
	info, err := captcha.GetCaptcha(ctx)
	if err != nil {
		return nil, err
	}
	return info, nil
}
