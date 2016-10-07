package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/pa001024/netdialer"
	"github.com/pa001024/netdialer/router"
	"github.com/pa001024/reflex/util"
)

const (
	VERSION = "netdialer CLI ver.0.4.2 by pa001024"
	EXAMPLE = "usage example: > netdialer -ip 10.0.10.20 -u 18123123122@ZZZ.XY -p 123456"
)

type Config struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RouterType string `json:"routerType"`
	RouterAddr string `json:"routerAddr"`
	RouterUser string `json:"routerUser"`
	RouterPwd  string `json:"routerPwd"`
}

var config *Config

func main() {
	username := flag.String("u", "", "shanxun username")
	password := flag.String("p", "", "shanxun password")
	isRouter := flag.Bool("r", false, "use router mode 2? (defalut off)")
	isDisconnect := flag.Bool("d", false, "disconnect?")
	realIP := flag.String("ip", "local", "real IP, 'stdin' or 'local' works router support [hiwifi openwrt asus]")
	raddress := flag.String("ra", "192.168.1.1", "router address")
	rusername := flag.String("ru", "root", "router username")
	rpassword := flag.String("rp", "admin", "router password")
	flag.Parse()
	config = &Config{*username, *password, *realIP, *raddress, *rusername, *rpassword}

	if !*isDisconnect {
		if config.Username == "" && !*isDisconnect {
			println(VERSION)
			println(EXAMPLE)
			flag.Usage()
			return
		} else if len(config.Username) < 12 || len(config.Password) != 6 {
			util.ERROR.Log("Invalid username or password")
			println(VERSION)
			println(EXAMPLE)
			flag.Usage()
			return
		}
	}
	d := netdialer.NewDialer(config.Username, config.Password)
	if config.RouterType != "" && config.RouterType != "local" {
		if config.RouterType == "stdin" {
			in := bufio.NewReader(os.Stdin)
			str, _ := in.ReadString(0)
			d.UserIP = str
		} else {
			d.UserIP = selectMode(config.RouterType)
		}
	}
	var err error
	if *isRouter {
		err = d.ConnectRouter()
	} else {
		if *isDisconnect {
			err = d.DisconnectDirect()
		} else {
			err = d.ConnectDirect()
		}
	}
	if err != nil {
		util.ERROR.Log(err)
	}
}
func selectMode(typ string) (rst string) {
	switch typ {
	case "hiwifi":
		rst = router.GetLanIP_HiwifiV2(config.RouterAddr, config.RouterPwd)
		if rst == "" {
			rst = router.GetLanIP_Hiwifi(config.RouterAddr, config.RouterPwd)
		}
	case "openwrt":
		rst = router.GetLanIP_Openwrt(config.RouterAddr, config.RouterPwd)
	case "asus":
		rst = router.GetLanIP_Asus(config.RouterAddr, config.RouterUser, config.RouterPwd)
	default:
		rst = config.RouterType
	}
	return
}
