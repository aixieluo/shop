package data_test

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"shop/app/user/internal/conf"
	"shop/app/user/internal/data"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var cleaner func()
var Db *data.Data
var ctx context.Context

func initialize(db *gorm.DB) error {
	err := db.AutoMigrate(&data.User{})
	return errors.WithStack(err)
}

var _ = BeforeSuite(func() {
	con, f := data.DockerMysql("mariadb", "latest")
	cleaner = f
	config := &conf.Data{Database: &conf.Data_Database{Source: con, Driver: "mysql"}}
	db := data.NewDB(config)
	MySqlDb, _, err := data.NewData(config, nil, db, nil)
	if err != nil {
		return
	}
	Db = MySqlDb
	err = initialize(db)
	if err != nil {
		return
	}
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	cleaner()
})
