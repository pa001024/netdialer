package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetLanIP_TPLink(address, username, password string) string {
	// TODO
	return ""
}
func SetWanInfo_TPLink(address, username, password, wanUser, wanPwd string) {
	req, err := http.NewRequest("GET", "http://"+address+"/userRpm/PPPoECfgRpm.htm?wan=0&wantype=2&acc="+
		strings.Replace(strings.Replace(url.QueryEscape(wanUser), "+", "%20", -1), "%40", "@", -1)+
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
func main() {
	address := flag.String("a", "192.168.1.1", "admin's address")
	username := flag.String("u", "admin", "admin's username")
	password := flag.String("p", "admin", "admin's password")
	isSetWan := flag.Bool("wan", false, "set wan info?")
	flag.Parse()
	if *isSetWan {
		lineIn := bufio.NewReader(os.Stdin)
		wanUser, _ := lineIn.ReadString('\n')
		wanPwd, _ := lineIn.ReadString('\n')
		os.Stdin.Close()
		SetWanInfo_TPLink(*address, *username, *password, wanUser, wanPwd)
	} else {
		ip := GetLanIP_TPLink(*address, *username, *password)
		fmt.Println(ip)
	}
}
