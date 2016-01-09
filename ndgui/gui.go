package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"github.com/pa001024/netdialer"
	"github.com/pa001024/netdialer/router"
	"github.com/pa001024/reflex/util"
)

type Config struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RouterAddr string `json:"routerAddr"`
	RouterUser string `json:"routerUser"`
	RouterPwd  string `json:"routerPwd"`
	RouterType string `json:"routerType"`
}

func main() {
	fout, _ := os.Create("dialer.log")
	defer fout.Close()
	bo := bufio.NewWriter(fout)
	defer bo.Flush()
	util.ERROR.SetOutput(bo)

	config := &Config{
		Username:   "18xxxxxxxxx@XXXX.XX",
		Password:   "******",
		RouterAddr: "192.168.1.1",
		RouterUser: "root",
		RouterPwd:  "admin",
		RouterType: "openwrt", // asus hiwifi openwrt
	}
	bin, err := ioutil.ReadFile("config.json")
	if err == nil {
		json.Unmarshal(bin, config)
	} else {
		util.ERROR.Log(err)
	}

	var usr, pwd *walk.LineEdit
	var raddr, rusr, rpwd *walk.LineEdit
	var lb, rb *walk.PushButton
	var mode *walk.ComboBox
	var mw *walk.MainWindow
	var db *walk.DataBinder
	go func() {
		for mw == nil {
			runtime.Gosched()
		}
		ic, err := walk.NewIconFromResourceId(6)
		if err == nil {
			// func onLoad() {
			db.SetAutoSubmit(true)
			mw.SetIcon(ic)
			// }
		}
	}()
	MainWindow{
		AssignTo: &mw,
		Title:    "闪讯拨号器GUI v0.1.1",
		MinSize:  Size{340, 0},
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:   &db,
			DataSource: config,
		},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Label{Text: "用户名", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("Username"), AssignTo: &usr, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "密码", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("Password"), AssignTo: &pwd, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "模式", MaxSize: Size{60, 20}},
					ComboBox{Editable: true, Value: Bind("RouterType"), Model: []string{"10.0.x.x(手动填写)", "openwrt", "hiwifi", "asus"}, AssignTo: &mode, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "路由地址", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("RouterAddr"), AssignTo: &raddr, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "路由用户名", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("RouterUser"), AssignTo: &rusr, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "路由密码", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("RouterPwd"), AssignTo: &rpwd, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					PushButton{
						AssignTo: &lb,
						Text:     "本地拨号",
						OnClicked: func() {
							d := netdialer.NewDialer(usr.Text(), pwd.Text())
							d.ConnectDirect()
						},
					},
					PushButton{
						AssignTo: &rb,
						Text:     "路由拨号",
						OnClicked: func() {
							d := netdialer.NewDialer(usr.Text(), pwd.Text())
							switch mode.Text() {
							case "hiwifi":
								d.UserIP = router.GetLanIP_Hiwifi(config.RouterAddr, config.RouterPwd)
							case "openwrt":
								d.UserIP = router.GetLanIP_Openwrt(config.RouterAddr, config.RouterPwd)
							case "asus":
								d.UserIP = router.GetLanIP_Asus(config.RouterAddr, config.RouterUser, config.RouterPwd)
							default:
								d.UserIP = config.RouterType
							}
							if d.UserIP == "" {
								d.RefreshIP()
							}
							d.ConnectDirect()
						},
					},
				}, MaxSize: Size{0, 20},
			},
		},
	}.Run()
	saveConfig(config)
}

func saveConfig(config *Config) {
	b, err := json.MarshalIndent(config, "", "\t")
	if err == nil {
		ioutil.WriteFile("config.json", b, 0644)
	} else {
		util.ERROR.Log(err)
	}
}
