package netdialer

import (
	"fmt"
	t "testing"

	// "github.com/pa001024/netdialer/router"
	// "github.com/pa001024/reflex/util"
)

var dialer *Dialer

func init() {
	dialer = NewDialer("17802201234@HYXY.XY", "123456")
}

func TestGetCryptUsername(t *t.T) {
	out := dialer.getCryptUsername()
	fmt.Println(out)
}

// func TestRouter(t *t.T) {
// 	out := router.GetLanIP_Hiwifi("192.168.199.1", "ai941024")
// 	fmt.Println(out)
// }

// func TestRouter2(t *t.T) {
// 	out := router.GetLanIP_Openwrt("192.168.99.1", "wewewewe")
// 	fmt.Println(out)
// }
