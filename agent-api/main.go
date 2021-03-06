package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-in-cn/XConf/agent-api/config"
	"github.com/micro-in-cn/XConf/agent-api/handler"
	pconfig "github.com/micro-in-cn/XConf/proto/config"
	"github.com/micro/cli"
	"github.com/micro/go-micro/web"
)

func main() {
	var cacheSize int

	service := web.NewService(
		web.Name("go.micro.api.agent"),
		web.Flags(
			cli.IntFlag{
				Name:        "cache_size",
				Usage:       "cache size (k)",
				EnvVar:      "XCONF_CACHE_SIZE",
				Value:       1024 * 1024,
				Destination: nil,
			},
		),
		web.Action(func(context *cli.Context) {
			cacheSize = context.Int("cache_size")
		}),
	)

	if err := service.Init(); err != nil {
		panic(err)
	}

	client := pconfig.NewConfigService("go.micro.srv.config", service.Options().Service.Client())

	config.Init(client, cacheSize*1024)
	router := Router()
	service.Handle("/", router)

	if err := service.Run(); err != nil {
		panic(err)
	}
}

func Router() *gin.Engine {
	router := gin.Default()
	r := router.Group("/agent/api/v1")
	r.GET("/config", handler.ReadConfig)
	r.GET("/config/raw", handler.ReadConfigRaw)
	r.GET("/watch", handler.WatchUpdate)
	r.GET("/watch/raw", handler.WatchUpdateRaw)

	return router
}
