package data

import (
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"shop/app/user/internal/conf"
	"time"

	slog "log"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewUserRepo)

// Data .
type Data struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, rdb *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, rdb: rdb}, cleanup, nil
}

func NewDB(c *conf.Data) *gorm.DB {
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags),
		logger.Config{
			SlowThreshold: 300 * time.Millisecond, // 慢查询阈值
			LogLevel:      logger.Info,
		},
	)
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Errorf("failed to opening connection to sqlite: %v", err)
		panic("failed to connect to database")
	}

	return db
}

func NewRedis(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Network:      c.Redis.Network,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	if err := rdb.Close(); err != nil {
		log.Error(err)
	}

	return rdb
}
