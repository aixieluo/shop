package server

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	netHttp "net/http"

	v1 "shop/api/shop/v1"
	"shop/internal/conf"
	"shop/internal/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, authService *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ResponseEncoder(func(w netHttp.ResponseWriter, r *netHttp.Request, i any) error {
			type response struct {
				Code    int    `json:"code"`
				Data    any    `json:"data"`
				Message string `json:"message"`
			}
			res := &response{
				Code:    netHttp.StatusOK,
				Data:    i,
				Message: "success",
			}
			msRes, err := json.Marshal(res)
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(msRes)
			return err
		}),
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	openAPIHandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIHandler)
	v1.RegisterAuthHTTPServer(srv, authService)
	return srv
}
