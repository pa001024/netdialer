package router

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"encoding/json"
)

func GetLanIP_TPLink(address, username, password string) string {
	// TODO
	return ""
}
func SetWanInfo_TPLink(address, username, password, wanUser, wanPwd string) {
	req, err := http.NewRequest("GET", "http://"+address+"/userRpm/PPPoECfgRpm.htm?wan=0&wantype=2&acc="+
		strings.Replace(strings.Replace(url.QueryEscape(url.QueryEscape(wanUser)), "+", "%20", -1), "%40", "@", -1)+
		"&psw="+wanPwd+"&confirm="+wanPwd+
		"&specialDial=0&SecType=0&sta_ip=0.0.0.0&sta_mask=0.0.0.0&linktype=4&waittime2=0&Connect=%C1%AC+%BD%D3", nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "text/plain; Charset=UTF-8")
	req.Header.Set("Connection", "Close")
	req.Header.Set("Referer", req.URL.String())

	req.SetBasicAuth(username, password)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	res.Body.Close()
	return
}

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

func GetLanIP_HiwifiV2(address, password string) string {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.Get("http://" + address + "/cgi-bin/turbo/admin_web/login_admin?" + url.Values{"username": {"admin"}, "password": {password}}.Encode())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	str := string(bin)
	ex := regexp.MustCompile(`"stok": "(.+?)",`) // /cgi-bin/turbo/;stok=8ab11a0e6a60d45e3658a4f81b0f2884
	li := ex.FindStringSubmatch(str)
	if len(li) > 1 {
		postBody := bytes.NewBufferString(`{"method":"network.wan.get_simple_info","data":{}}`)
		res, err := client.Post("http://"+address+"/cgi-bin/turbo"+li[1]+"/proxy/call", "application/x-www-form-urlencoded", postBody)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		bin, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		str = string(bin)
		ex = regexp.MustCompile(`"wan_ip":"(10\..+?)"`)
		li = ex.FindStringSubmatch(str)
		if len(li) > 1 {
			return li[1]
		}
	}
	return ""
}

func GetLanIP_HiwifiV3(address, password string) string {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.Get("http://" + address + "/cgi-bin/turbo/api/login/login_admin?" + url.Values{"username": {"admin"}, "password": {password}}.Encode())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	str := string(bin)
	ex := regexp.MustCompile(`"/;stok=(.+?)",`) // /cgi-bin/turbo/;stok=8ab11a0e6a60d45e3658a4f81b0f2884
	li := ex.FindStringSubmatch(str)
	if len(li) > 1 {
		postBody := bytes.NewBufferString(`{"method":"wan.get_status","data":{},"lang":"zh-CN","version":"v1"}`)
		res, err := client.Post("http://"+address+"/cgi-bin/turbo/;stok="+li[1]+"/proxy/call", "application/json", postBody)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		bin, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		str = string(bin)
		ex = regexp.MustCompile(`"wan_ip":"(10\..+?)"`)
		li = ex.FindStringSubmatch(str)
		if len(li) > 1 {
			return li[1]
		}
	}
	return ""
}

func SetWanInfo_HiwifiV3(address, password, wanUser, wanPwd string) (err error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.Get("http://" + address + "/cgi-bin/turbo/api/login/login_admin?" + url.Values{"username": {"admin"}, "password": {password}}.Encode())
	if err != nil {
		fmt.Println(err)
		return
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	str := string(bin)
	ex := regexp.MustCompile(`"/;stok=(.+?)",`) // /cgi-bin/turbo/;stok=8ab11a0e6a60d45e3658a4f81b0f2884
	li := ex.FindStringSubmatch(str)
	if len(li) > 1 {
		bin, _ = json.Marshal(map[string]interface{}{
			"method": "wan.set_pppoe_way",
			"data": map[string]string{
				"type":              "pppoe",
				"peerdns":           "0",
				"pppoe_name":        wanUser,
				"pppoe_passwd":      wanPwd,
				"switch_status_wan": "auto",
				"special_dial":      "0",
				"override_dns":      "119.29.29.29",
				"override_dns2":     "180.76.76.76",
				"mac_clone":         "0",
			},
			"lang":    "zh-CN",
			"version": "v1",
		})
		fmt.Println(wanUser, string(bin))
		postBody := bytes.NewBuffer(bin)
		res, err = client.Post("http://"+address+"/cgi-bin/turbo/;stok="+li[1]+"/proxy/call", "application/json", postBody)
		if err != nil {
			fmt.Println(err)
			return
		}
		bin, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		str = string(bin)
		ex = regexp.MustCompile(`成功`)
		li = ex.FindStringSubmatch(str)
		if len(li) > 1 {
			return nil
		}
	}
	return
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
