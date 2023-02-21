package testdata

import "shop/app/user/internal/biz"

func User() *biz.User {
	return &biz.User{
		ID:       1,
		Username: "aixieluo",
		Nickname: "39",
		Email:    "aixieluo@gmail.com",
		Mobile:   "17521066239",
		Password: "123123",
	}
}
