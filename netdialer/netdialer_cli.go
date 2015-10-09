package main

import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/pa001024/netdialer"
	"github.com/pa001024/reflex/util"
)

func main() {
	username := flag.String("u", "", "shanxun username")
	password := flag.String("p", "", "shanxun password")
	isRouter := flag.Bool("r", false, "use router?")
	routerIP := flag.String("ra", "192.168.1.1", "router IP")
	routerUser := flag.String("ru", "admin", "router User")
	routerPwd := flag.String("rp", "admin", "router password")
	noCheck := flag.Bool("nc", false, "no check:dont need check?")
	realIP := flag.String("ip", "local", "real IP, 'stdin' or 'local' works")
	flag.Parse()
	if len(*username) < 12 || len(*password) != 6 {
		util.ERROR.Log("Invalid username or password")
		flag.Usage()
		return
	}
	d := netdialer.NewDialer(*username, *password)
	if *realIP != "" && *realIP != "local" {
		if *realIP == "stdin" {
			in := bufio.NewReader(os.Stdin)
			str, _ := in.ReadString(0)
			d.UserIP = str
		} else {
			d.UserIP = *realIP
		}
	}
	var err error
	if *isRouter {
		d.Router.Addr = *routerIP
		d.Router.User = *routerUser
		d.Router.Pwd = *routerPwd
		err = d.ConnectRouter()
		if !*noCheck {
			check()
		}
	} else {
		err = d.ConnectDirect()
		if !*noCheck {
			check()
		}
	}
	if err != nil {
		util.ERROR.Log(err)
	}
}

func check() {
	<-time.After(time.Second * 5)
	ip := util.GetIPInfo()
	if ip != nil && ip.IP != "" {
		util.INFO.Log("连接成功, 当前IP:", ip.IP)
	} else {
		util.INFO.Log("连接失败")
	}
}
