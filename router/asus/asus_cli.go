package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetLanIP_Asus() string {
	address := flag.String("a", "192.168.1.1", "admin's address")
	username := flag.String("u", "admin", "admin's username")
	password := flag.String("p", "admin", "admin's password")
	flag.Parse()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+*address+"/status.asp", nil)
	req.SetBasicAuth(*username, *password)
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
func main() {
	ip := GetLanIP_Asus()
	fmt.Println(ip)
}
