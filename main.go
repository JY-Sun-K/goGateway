package main

import (
	"flag"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"go_gateway/dao"
	"go_gateway/http_proxy_router"
	"go_gateway/router"
	"os"
	"os/signal"
	"syscall"
)

var (
	endpoint =flag.String("endpoint","","input endpoint dashboard or server")
	config =flag.String("config","","input config file like ./conf/dev/")
)

func main()  {
	flag.Parse()
	if *endpoint=="" {
		flag.Usage()
		os.Exit(1)
	}
	if *config=="" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint=="dashboard" {
		lib.InitModule(*config,[]string{"base","mysql","redis"})
		defer lib.Destroy()


		router.HttpServerRun()



		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	}else if *endpoint=="server" {
		lib.InitModule(*config,[]string{"base","mysql","redis"})
		defer lib.Destroy()
		err:=dao.ServiceManagerHandler.LoadOnce()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go func() {
			http_proxy_router.HttpServerRun()
		}()
		go func() {
			http_proxy_router.HttpsServerRun()
		}()
		fmt.Println("run server")
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit



		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}else {
		os.Exit(1)
	}
	
	

}