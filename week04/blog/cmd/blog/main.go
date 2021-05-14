package main

import (
	"context"
	"flag"
	"gopkg.in/yaml.v2"
	"my"
	"os"
	"time"

	"blog/internal/conf"
	"my/config"
	"my/config/file"
	"my/log"
	"my/transport/grpc"
	"my/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *my.App {
	return my.New(
		my.Name(Name),
		my.Version(Version),
		my.Metadata(map[string]string{}),
		my.Logger(logger),
		my.Server(
			hs,
			gs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)

	config := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := config.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := config.Scan(&bc); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
	}(ctx)

	app, cleanup, err := initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
