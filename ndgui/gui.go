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

var (
	config *Config = &Config{
		Username:   "18xxxxxxxxx@XXXX.XX",
		Password:   "******",
		RouterAddr: "192.168.1.1",
		RouterUser: "root",
		RouterPwd:  "admin",
		RouterType: "local", // asus hiwifi openwrt
	}
)

func main() {
	fout, _ := os.Create("dialer.log")
	defer fout.Close()
	bo := bufio.NewWriter(fout)
	defer bo.Flush()
	util.ERROR.SetOutput(bo)

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
	var laddr, lusr, lpwd *walk.Splitter
	go func() {
		for mw == nil {
			runtime.Gosched()
		}
		ic, err := walk.NewIconFromResourceId(6)
		if err == nil {
			// func onLoad() {
			db.SetAutoSubmit(true)
			mw.SetIcon(ic)
			switch config.RouterType {
			case "hiwifi":
				laddr.SetVisible(true)
				lusr.SetVisible(false)
				lpwd.SetVisible(true)
			case "openwrt":
				laddr.SetVisible(true)
				lusr.SetVisible(true)
				lpwd.SetVisible(true)
			case "asus":
				laddr.SetVisible(true)
				lusr.SetVisible(true)
				lpwd.SetVisible(true)
			default:
				laddr.SetVisible(false)
				lusr.SetVisible(false)
				lpwd.SetVisible(false)
			}
			mw.SetSize(walk.Size{0, 0})
			// }
		}
	}()
	MainWindow{
		AssignTo: &mw,
		Title:    "闪讯拨号器GUI v0.4.1",
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
					LineEdit{Text: Bind("Password"), AssignTo: &pwd, MaxSize: Size{0, 20}, PasswordMode: true,
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							pwd.SetPasswordMode(false)
						},
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							pwd.SetPasswordMode(true)
						},
					},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "模式", MaxSize: Size{60, 20}},
					ComboBox{AssignTo: &mode,
						Editable: true, Value: Bind("RouterType"),
						Model:   []string{"local", "10.0.x.x(手动填写)", "openwrt", "hiwifi", "asus"},
						MaxSize: Size{0, 20},
						OnCurrentIndexChanged: func() {
							switch mode.CurrentIndex() {
							case 2: //"openwrt":
								laddr.SetVisible(true)
								lusr.SetVisible(true)
								lpwd.SetVisible(true)
							case 3: //"hiwifi":
								config.RouterAddr = "192.168.199.1"
								laddr.SetVisible(true)
								lusr.SetVisible(false)
								lpwd.SetVisible(true)
							case 4: //"asus":
								laddr.SetVisible(true)
								lusr.SetVisible(true)
								lpwd.SetVisible(true)
							default:
								laddr.SetVisible(false)
								lusr.SetVisible(false)
								lpwd.SetVisible(false)
							}
							mw.SetSize(walk.Size{0, 0})
						},
					},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				AssignTo: &laddr,
				Children: []Widget{
					Label{Text: "路由地址", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("RouterAddr"), AssignTo: &raddr, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				AssignTo: &lusr,
				Children: []Widget{
					Label{Text: "路由用户名", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("RouterUser"), AssignTo: &rusr, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				AssignTo: &lpwd,
				Children: []Widget{
					Label{Text: "路由密码", MaxSize: Size{60, 20}},
					LineEdit{Text: Bind("RouterPwd"), AssignTo: &rpwd, MaxSize: Size{0, 20}, PasswordMode: true,
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							rpwd.SetPasswordMode(false)
						},
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							rpwd.SetPasswordMode(true)
						},
					},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					PushButton{
						AssignTo: &lb,
						Text:     "开始连接",
						OnClicked: func() {
							if mode.Text() == "10.0.x.x(手动填写)" {
								walk.MsgBox(mw, "请填写IP", "手动填写需要自己获取IP 你可在路由器中自己查找", walk.MsgBoxOK)
								return
							}
							lb.SetText("连接中...")
							lb.SetEnabled(false)
							rb.SetEnabled(false)
							go func() {
								d := netdialer.NewDialer(usr.Text(), pwd.Text())
								d.UserIP = selectMode(mode.Text())
								d.ConnectDirect()
								d = nil
								lb.SetEnabled(true)
								rb.SetEnabled(true)
								lb.SetText("开始连接")
								if err == nil {
									walk.MsgBox(mw, "连接成功", "感谢使用", walk.MsgBoxOK)
								} else {
									walk.MsgBox(mw, "连接失败", err.Error(), walk.MsgBoxOK)
								}
							}()
						},
					},
					PushButton{
						AssignTo: &rb,
						Text:     "断开连接",
						OnClicked: func() {
							if mode.Text() == "10.0.x.x(手动填写)" {
								walk.MsgBox(mw, "请填写IP", "手动填写需要自己获取IP 你可在路由器中自己查找", walk.MsgBoxOK)
								return
							}
							rb.SetText("断开中...")
							lb.SetEnabled(false)
							rb.SetEnabled(false)
							go func() {
								d := netdialer.NewDialer(usr.Text(), pwd.Text())
								d.UserIP = selectMode(mode.Text())
								err := d.DisconnectDirect()
								d = nil
								lb.SetEnabled(true)
								rb.SetEnabled(true)
								rb.SetText("断开连接")
								if err == nil {
									walk.MsgBox(mw, "断开成功", "感谢使用", walk.MsgBoxOK)
								} else {
									walk.MsgBox(mw, "断开失败", err.Error(), walk.MsgBoxOK)
								}
							}()
						},
					},
				}, MaxSize: Size{0, 20},
			},
		},
	}.Run()
	saveConfig(config)
}

func selectMode(typ string) (rst string) {
	switch typ {
	case "hiwifi":
		rst = router.GetLanIP_Hiwifi(config.RouterAddr, config.RouterPwd)
	case "openwrt":
		rst = router.GetLanIP_Openwrt(config.RouterAddr, config.RouterPwd)
	case "asus":
		rst = router.GetLanIP_Asus(config.RouterAddr, config.RouterUser, config.RouterPwd)
	default:
		rst = config.RouterType
	}
	return
}

func saveConfig(config *Config) {
	b, err := json.MarshalIndent(config, "", "\t")
	if err == nil {
		ioutil.WriteFile("config.json", b, 0644)
	} else {
		util.ERROR.Log(err)
	}
}
