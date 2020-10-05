package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/northfun/house/src/app"
)

var (
	config string
)

func main() {
	flag.StringVar(&config, "c", "", "config path")
	flag.Parse()

	var app app.App
	app.Init(config)

	app.Start()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGKILL)
	s := <-c

	fmt.Println("get signal:", s)

	app.Stop()
}
