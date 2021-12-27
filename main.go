package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/web"
	"github.com/taadis/blog-web/conf"
	"github.com/taadis/blog-web/internal/pkg/mysql"
	"github.com/taadis/blog-web/internal/routers"
	logger "github.com/taadis/blog-web/pkg/log"
)

var (
	cfgFile string
)

func beforeStart() error {
	log.Printf("before start...init db,redis,es,etc.")

	err := config.Load(file.NewSource(file.WithPath("./conf/dev.yml")))
	if err != nil {
		log.Printf("beforeStart config.LoadFile error:%+v", err)
		return err
	}
	var cfg conf.Config
	err = config.Scan(&cfg)
	if err != nil {
		log.Printf("beforeStart config.Scan error:%+v", err)
		return err
	}

	// init config
	//cfg := conf.Init(cfgFile)

	// init logger
	logger.Init(&cfg.Logger)

	// init redis
	//redis.Init(&cfg.Redis)

	// init orm
	//model.Init(&cfg.ORM)

	// init mysql
	mysql.Init(&cfg.Mysql)
	return nil
}

func main() {
	service := web.NewService(
		web.Name("go.micro.web.blog"),
		web.Registry(etcd.NewRegistry(registry.Addrs(os.Getenv("MICRO_REGISTRY_ADDRESS")))),
		web.BeforeStart(beforeStart),
	)
	service.Init()
	//service.HandleFunc("/", index)
	service.Handle("/", routers.InitRouter())
	service.Run()
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `web`)
}
