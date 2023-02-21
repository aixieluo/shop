package data

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"

	"shop/app/user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID        uint       `gorm:"primarykey"`
	Username  string     `json:"username,omitempty" gorm:"index:idx_username;unique;type:varchar(20) comment '用户名';not null"`
	Nickname  string     `json:"nickname,omitempty" gorm:"idx_nickname;type:varchar(40) comment '昵称';not null"`
	Email     string     `json:"email,omitempty" gorm:"idx_email;type:varchar(50) comment '邮箱';not null"`
	Mobile    string     `json:"mobile,omitempty" gorm:"idx_mobile;type:varchar(11) comment '手机号';not null"`
	Password  string     `json:"password,omitempty" gorm:"type:varchar(100) comment '密码';not null"`
	Gender    uint8      `json:"gender,omitempty" gorm:"default: 0;type:tinyint(1) comment '性别0未知1男2女';not null"`
	Birthday  *time.Time `json:"birthday,omitempty" gorm:"type:datetime"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func (r *UserRepo) List(ctx context.Context, page int, perPage int) ([]*biz.User, int, error) {
	var users []User
	var total int64
	r.data.db.Model(&User{}).Count(&total)
	tx := r.data.db.Scopes(paginate(page, perPage)).Find(&users)
	if tx.Error != nil {
		return nil, 0, status.Errorf(codes.Internal, tx.Error.Error())
	}
	rv := make([]*biz.User, 0)
	for _, user := range users {
		rv = append(rv, &biz.User{
			ID:       uint64(user.ID),
			Nickname: user.Username,
			Email:    user.Nickname,
			Mobile:   user.Email,
		})
	}
	return rv, int(total), nil
}

func (r *UserRepo) FindByID(ctx context.Context, id uint) (*biz.User, error) {
	var user User
	err := first(ctx, r.data.db, &user, id)
	if err != nil {
		return nil, err
	}
	rsp := modelToResponse(user)
	return &rsp, nil
}

func (r *UserRepo) Create(ctx context.Context, bu *biz.User) (*biz.User, error) {
	var u User
	result := r.data.db.Where(&biz.User{Username: bu.Username}).First(&u)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	u.Username = bu.Username
	u.Nickname = bu.Nickname
	u.Password = encrypt(bu.Password)
	u.Mobile = bu.Mobile
	u.Email = bu.Email

	tx := r.data.db.Create(&u)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, tx.Error.Error())
	}
	rsp := modelToResponse(u)
	return &rsp, nil
}

func (r *UserRepo) Update(ctx context.Context, id uint, user *biz.User) (bool, error) {
	var u User
	err := first(ctx, r.data.db, &u, id)
	if err != nil {
		return false, err
	}

	u.Nickname = user.Nickname
	u.Email = user.Email
	u.Mobile = user.Mobile
	tx := r.data.db.Save(&u)
	if tx.Error != nil {
		return false, status.Errorf(codes.Internal, tx.Error.Error())
	}

	return true, nil
}
