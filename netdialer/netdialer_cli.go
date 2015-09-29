package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/pa001024/netdialer"
	"github.com/pa001024/reflex/util"
)

func main() {
	username := flag.String("u", "", "shanxun username")
	password := flag.String("p", "", "shanxun password")
	isRouter := flag.Bool("r", true, "use router?")
	routerIP := flag.String("ra", "192.168.1.1", "router IP")
	routerUser := flag.String("ru", "admin", "router User")
	routerPwd := flag.String("rp", "admin", "router password")
	flag.Parse()
	d := netdialer.NewDialer(*username, *password)
	var err error
	if *isRouter {
		d.Router.Addr = *routerIP
		d.Router.User = *routerUser
		d.Router.Pwd = *routerPwd
		err = d.ConnectRouter()
		check()
	} else {
		err = d.ConnectDirect()
		check()
	}
	if err != nil {
		fmt.Errorf("%s", err)
	}
}

func check() {
	go func() {
		<-time.After(time.Second * 2)
		ip := util.GetIPInfo()
		if ip != nil {
			util.INFO.Log("连接成功, 当前IP:", ip.IP)
		}
	}()
}
