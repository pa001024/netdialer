package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

func GetLanIP_Asus(address, username, password string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+address+"/status.asp", nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	str := string(bin)
	ex := regexp.MustCompile(`wanlink_ipaddr\(\) \{ return '(.+?)';\}`)
	li := ex.FindStringSubmatch(str)
	if len(li) > 1 {
		return li[1]
	}
	return ""
}

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

func GetLanIP_Openwrt(address, password string) string {
	// Login first
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.PostForm("http://"+address+"/", url.Values{"luci_username": {"root"}, "luci_password": {password}})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	str := string(bin)
	ex := regexp.MustCompile(`/cgi-bin/luci/;stok=([a-z0-9]{32})`) // /cgi-bin/luci/;stok=dfc41c0ba4035a36922a6df4e26f6dd7/
	li := ex.FindStringSubmatch(str)
	if len(li) == 0 {
		// Login ...
		jar, _ = cookiejar.New(nil)
		client = &http.Client{Jar: jar}
		res, err = client.PostForm("http://"+address+"/", url.Values{"username": {"root"}, "password": {password}})
		if err != nil {
			fmt.Println(err)
			return ""
		}
		bin, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		str = string(bin)
		ex = regexp.MustCompile(`/cgi-bin/luci/;stok=([a-z0-9]{32})`) // /cgi-bin/luci/;stok=dfc41c0ba4035a36922a6df4e26f6dd7/
		li = ex.FindStringSubmatch(str)
	}
	if len(li) > 1 {
		res, err = client.Get("http://" + address + li[0] + "?status=1")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		bin, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		str = string(bin)
		ex = regexp.MustCompile(`"ipaddr":\s*"(10\.[\.0-9]+?)",`)
		li = ex.FindStringSubmatch(str)
		if len(li) > 1 {
			return li[1]
		}
	}
	return ""
}
