package data

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"shop/app/user/internal/biz"
	"strings"
)

func encrypt(pwd string) string {
	options := &password.Options{SaltLen: 16, Iterations: 10000, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(pwd, options)
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}

func verify(pwd, encryptedPwd string) (bool, error) {
	options := &password.Options{SaltLen: 16, Iterations: 10000, KeyLen: 32, HashFunction: sha512.New}
	passwordSlice := strings.Split(encryptedPwd, "$")
	check := password.Verify(pwd, passwordSlice[2], passwordSlice[3], options)
	return check, nil
}

func modelToResponse(user User) biz.User {
	return biz.User{
		ID:       uint64(user.ID),
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
	}
}

func first(ctx context.Context, db *gorm.DB, dest any, id uint) error {
	tx := db.First(dest, id)
	if tx.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "用户不存在")
	}
	return nil
}

func paginate(page, perPage int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case perPage > 100:
			perPage = 100
		case perPage <= 0:
			perPage = 10
		}
		offset := (page - 1) * perPage
		return db.Offset(offset).Limit(perPage)
	}
}
