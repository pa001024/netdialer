package main

import (
	// "bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	// "os"
	"regexp"
)

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
	if len(li) > 1 {
		res, err = client.Get("http://" + address + li[0] + "?status=1")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		bin, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()
		str = string(bin)
		ex = regexp.MustCompile(`"ipaddr":"(10\.[\.0-9]+?)",`)
		li = ex.FindStringSubmatch(str)
		if len(li) > 1 {
			return li[1]
		}
	}
	return ""
}
func SetWanInfo_Openwrt(address, password, wanUser, wanPwd string) {
	/// todo
	return
}
func main() {
	address := flag.String("a", "192.168.99.1", "admin's address")
	password := flag.String("p", "admin", "admin's password")
	// isSetWan := flag.Bool("wan", false, "set wan info?")
	flag.Parse()
	if false {
		// lineIn := bufio.NewReader(os.Stdin)
		// wanUser, _ := lineIn.ReadString('\n')
		// wanPwd, _ := lineIn.ReadString('\n')
		// wanUser, _ = url.QueryUnescape(wanUser)
		// os.Stdin.Close()
		// SetWanInfo_Hiwifi(*address, *password, wanUser, wanPwd)
	} else {
		ip := GetLanIP_Openwrt(*address, *password)
		fmt.Println(ip)
	}
}
