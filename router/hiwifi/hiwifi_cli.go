package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
)

func GetLanIP_Hiwifi(address, password string) string {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.PostForm("http://"+address+"/cgi-bin/turbo/admin_web", url.Values{"username": {"admin"}, "password": {password}})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	str := string(bin)
	ex := regexp.MustCompile(`URL_ROOT_PATH = "(.+?)";`) // /cgi-bin/turbo/;stok=8ab11a0e6a60d45e3658a4f81b0f2884
	li := ex.FindStringSubmatch(str)
	if len(li) > 1 {
		res, err = client.Get("http://" + address + "" + li[1] + "/api/network/get_wan_info")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		bin, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		str = string(bin)
		ex = regexp.MustCompile(`"ip": "(10\..+?)" }`)
		li = ex.FindStringSubmatch(str)
		if len(li) > 1 {
			return li[1]
		}
	}
	return ""
}
func SetWanInfo_Hiwifi(address, password, wanUser, wanPwd string) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.PostForm("http://"+address+"/cgi-bin/turbo/admin_web", url.Values{"username": {"admin"}, "password": {password}})
	if err != nil {
		fmt.Println(err)
		return
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	str := string(bin)
	ex := regexp.MustCompile(`URL_ROOT_PATH = "(.+?)";`) // /cgi-bin/turbo/;stok=8ab11a0e6a60d45e3658a4f81b0f2884
	li := ex.FindStringSubmatch(str)
	if len(li) > 1 {
		res, err = client.PostForm("http://"+address+""+li[1]+"/api/network/set_wan_connect", url.Values{
			"network_type":      {"pppoe"},
			"type":              {"pppoe"},
			"pppoe_name":        {wanUser},
			"pppoe_passwd":      {wanPwd},
			"ip_type":           {"dhcp"},
			"static_ip":         {""},
			"static_mask":       {""},
			"static_gw":         {""},
			"static_dns":        {""},
			"static_dns2":       {""},
			"ssid":              {""},
			"channel":           {""},
			"bssid":             {""},
			"encryption":        {""},
			"ssid_select_mode":  {""},
			"key":               {""},
			"key_show":          {""},
			"peerdns":           {"0"},
			"override_dns":      {"119.29.29.29"},
			"override_dns2":     {"180.76.76.76"},
			"uptime":            {"0"},
			"switch_status_wan": {"auto"},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		res.Body.Close()
	}
	return
}
func main() {
	address := flag.String("a", "192.168.199.1", "admin's address")
	password := flag.String("p", "admin", "admin's password")
	isSetWan := flag.Bool("wan", false, "set wan info?")
	flag.Parse()
	if *isSetWan {
		lineIn := bufio.NewReader(os.Stdin)
		wanUser, _ := lineIn.ReadString('\n')
		wanPwd, _ := lineIn.ReadString('\n')
		wanUser, _ = url.QueryUnescape(wanUser)
		os.Stdin.Close()
		SetWanInfo_Hiwifi(*address, *password, wanUser, wanPwd)
	} else {
		ip := GetLanIP_Hiwifi(*address, *password)
		fmt.Println(ip)
	}
}
