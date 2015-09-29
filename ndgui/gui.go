package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"github.com/pa001024/netdialer"
	"github.com/pa001024/reflex/util"
)

type Config struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RouterAddr string `json:"routerAddr"`
	RouterUser string `json:"routerUser"`
	RouterPwd  string `json:"routerPwd"`
}

func main() {
	fout, _ := os.Create("dialer.log")
	defer fout.Close()
	bo := bufio.NewWriter(fout)
	defer bo.Flush()
	util.ERROR.SetOutput(bo)

	config := &Config{
		RouterAddr: "192.168.1.1",
		RouterUser: "admin",
		RouterPwd:  "admin",
	}
	bin, err := ioutil.ReadFile("config.json")
	if err == nil {
		json.Unmarshal(bin, config)
	} else {
		util.ERROR.Log(err)
	}

	var usr, pwd *walk.LineEdit
	MainWindow{
		Title:   "闪讯拨号器GUI",
		MinSize: Size{340, 0},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Label{Text: "用户名", MaxSize: Size{60, 20}},
					LineEdit{Text: config.Username, AssignTo: &usr, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "密码", MaxSize: Size{60, 20}},
					LineEdit{Text: config.Password, AssignTo: &pwd, MaxSize: Size{0, 20}},
				}, MaxSize: Size{0, 20},
			},
			HSplitter{
				Children: []Widget{
					PushButton{
						Text: "本地拨号",
						OnClicked: func() {
							config.Username = usr.Text()
							config.Password = pwd.Text()
							d := netdialer.NewDialer(usr.Text(), pwd.Text())
							d.ConnectDirect()
						},
					},
					PushButton{
						Text: "路由拨号",
						OnClicked: func() {
							config.Username = usr.Text()
							config.Password = pwd.Text()
							d := netdialer.NewDialer(usr.Text(), pwd.Text())
							d.ConnectRouter()
						},
					},
				}, MaxSize: Size{0, 20},
			},
		},
	}.Run()

	b, err := json.Marshal(config)
	if err == nil {
		ioutil.WriteFile("config.json", b, 0644)
	} else {
		util.ERROR.Log(err)
	}
}
