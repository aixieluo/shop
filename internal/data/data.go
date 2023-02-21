package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"time"

	userV1 "shop/api/user/v1"
	"shop/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAuthRepo, NewUserServiceClient, NewDiscovery)

// Data .
type Data struct {
	uc  userV1.UserClient
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, uc userV1.UserClient, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data"))
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		uc:  uc,
		log: helper,
	}, cleanup, nil
}

func NewUserServiceClient(ac *conf.Auth, sr *conf.Service, rr registry.Discovery) userV1.UserClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.User.Endpoint),
		grpc.WithTimeout(2*time.Second),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithDiscovery(rr),
	)
	if err != nil {
		panic(err)
	}
	c := userV1.NewUserClient(conn)
	return c
}

// NewRegistrar add consul
func NewRegistrar(conf *conf.Registry) registry.Registrar {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
