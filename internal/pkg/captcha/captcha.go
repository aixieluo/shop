package captcha

import (
	"context"
	"github.com/mojocn/base64Captcha"
)

var Store = base64Captcha.DefaultMemStore

type Info struct {
	CaptchaId string
	PicPath   string
}

func GetCaptcha(ctx context.Context) (*Info, error) {
	driver := base64Captcha.NewDriverDigit(80, 250, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, Store)
	id, b64s, err := cp.Generate()
	if err != nil {
		return nil, err
	}
	return &Info{
		CaptchaId: id,
		PicPath:   b64s,
	}, nil
}
