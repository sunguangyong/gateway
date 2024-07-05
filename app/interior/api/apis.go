package main

import (
	"flag"
	"fmt"
	"log"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/protocols"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/config"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/handler"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/apis.yaml", "the config file")

func init() {
	log.SetFlags(log.Llongfile | log.Ltime)
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	protocols.StartProtocols(ctx)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
