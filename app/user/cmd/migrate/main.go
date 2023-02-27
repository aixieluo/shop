package main

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"os"
	"shop/app/user/internal/conf"
	"shop/app/user/internal/data"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", fmt.Sprintf("%v/src/shop/app/user/configs", os.Getenv("GOPATH")), "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	db := data.NewDB(bc.Data)

	if err := db.AutoMigrate(&data.User{}); err != nil {
		panic(err)
	}
}
