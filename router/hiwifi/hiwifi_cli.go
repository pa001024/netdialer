package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

func GetLanIP_Hiwifi() string {
	password := flag.String("p", "", "admin's password")
	flag.Parse()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.PostForm("http://192.168.199.1/cgi-bin/turbo/admin_web", url.Values{"username": {"admin"}, "password": {*password}})
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
		res, err = client.Get("http://192.168.199.1" + li[1] + "/api/network/get_wan_info")
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
func main() {
	ip := GetLanIP_Hiwifi()
	fmt.Println(ip)
}
