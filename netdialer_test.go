package netdialer

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"fmt"
	t "testing"

	// "github.com/pa001024/reflex/util"
)

var dialer *Dialer

func init() {
	dialer = NewDialer("17805805321@HYXY.XY", "266174")
}

type TEMPSTR struct {
	UserIP   string `xml:"Redirect>UserIP"`
	LoginURL string `xml:"Redirect>LoginURL"`
	Uuid     string `xml:"Redirect>Uuid"`
}

func TestStep1(t *t.T) {
	par := url.Values{"wlanuserip": {"10.0.1.2"}}
	body := ioutil.NopCloser(strings.NewReader(par.Encode()))
	req, _ := http.NewRequest("POST", "http://115.239.134.163:8080/showlogin.do", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; ")
	req.Header.Set("User-Agent", "China Telecom Client")
	res, _ := http.DefaultClient.Do(req)
	de := xml.NewDecoder(res.Body)
	v := &TEMPSTR{}
	de.Decode(v)
	fmt.Println(v)
}

func TestGetCryptUsername(t *t.T) {
	out := dialer.getCryptUsername()
	fmt.Println(out)
}
