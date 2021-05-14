// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"blog/internal/biz"
	"blog/internal/conf"
	"blog/internal/data"
	"blog/internal/server"
	"blog/internal/service"
	"github.com/google/wire"
	"my"
	"my/log"
)

// initApp init application.
func initApp(*conf.Server, *conf.Data, log.Logger) (*my.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
