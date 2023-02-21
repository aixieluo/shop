package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "shop/app/user/api/user/v1"
)

var userClient v1.UserClient
var conn *grpc.ClientConn

func main() {
	Init()

	TestCreateUser()

	defer conn.Close()
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("grpc link err" + err.Error())
	}
	userClient = v1.NewUserClient(conn)
}

func TestCreateUser() {
	rsp, err := userClient.CreateUser(context.Background(), &v1.CreateUserRequest{
		Username: "aixiela",
		Nickname: "ajks",
		Email:    "adsjk@ask.com",
		Mobile:   "14141241241",
		Password: "123123",
	})
	if err != nil {
		panic("grpc 创建用户失败" + err.Error())
	}
	fmt.Println(rsp.Id)
}
