package data

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"log"
	"time"
)

func DockerMysql(img string, version string) (string, func()) {
	return innerDockerMysql(img, version)
}

func innerDockerMysql(img string, version string) (string, func()) {
	pool, err := dockertest.NewPool("")
	pool.MaxWait = 2 * time.Minute
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.Run(img, version, []string{"MYSQL_ROOT_PASSWORD=123123", "MYSQL_ROOT_HOST=%"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	confStr := fmt.Sprintf("root:123123@(localhost:%s)/?charset=utf8mb4&parseTime=True&loc=Local", resource.GetPort("3306/tcp"))
	confStr2 := fmt.Sprintf("root:123123@(localhost:%s)/shop?charset=utf8mb4&parseTime=True&loc=Local", resource.GetPort("3306/tcp"))

	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("mysql", confStr)
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
		err = db.Ping()
		if err != nil {
			return err
		}
		_, err = db.Exec("CREATE DATABASE shop")
		if err != nil {
			return err
		}
		_, err = db.Exec("USE shop")
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker : %s", err)
	}

	return confStr2, func() {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
