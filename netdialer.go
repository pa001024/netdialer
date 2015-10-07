package netdialer

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pa001024/reflex/util"
)

const (
	NET_PWD_MASK = "7%ChIna3#@Net*%"
	RADIUS       = "singlenet01" //"\x73\x69\x6e\x67\x6c\x65\x6e\x65\x74\x30\x31"
)

type Dialer struct {
	username    string
	rawPassword string
	password    string
	UserIP      string
	uuid        string
	logoffURL   string
	ratingtype  string
	Router      RouterInfo
}

type RouterInfo struct {
	Type int
	Addr string
	User string
	Pwd  string
}

const (
	Router_TP = iota
)

func NewDialer(username, password string) (obj *Dialer) {
	obj = &Dialer{
		username:   username,
		ratingtype: "1",
		Router: RouterInfo{
			Type: Router_TP,
			Addr: "192.168.1.1",
			User: "admin",
			Pwd:  "admin",
		},
	}
	obj.SetPassword(password)
	obj.RefreshIP()
	return
}

// 路由器拨号
func (this *Dialer) ConnectRouter() (err error) {
	defer util.Catch(&err)
	req, err := http.NewRequest("GET", "http://"+this.Router.Addr+"/userRpm/PPPoECfgRpm.htm?wan=0&wantype=2&acc="+
		strings.Replace(strings.Replace(url.QueryEscape(this.getCryptUsername()), "+", "%20", -1), "%40", "@", -1)+
		"&psw="+this.rawPassword+"&confirm="+this.rawPassword+
		"&specialDial=0&SecType=0&sta_ip=0.0.0.0&sta_mask=0.0.0.0&linktype=4&waittime2=0&Connect=%C1%AC+%BD%D3", nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "text/plain; Charset=UTF-8")
	req.Header.Set("Connection", "Close")
	req.Header.Set("Referer", req.URL.String())

	req.SetBasicAuth(this.Router.User, this.Router.Pwd)
	util.Try(err)
	res, err := http.DefaultClient.Do(req)
	util.Try(err)
	res.Body.Close()
	return
}

//本地拨号
func (this *Dialer) ConnectDirect() (err error) {
	defer util.Catch(&err)
	info, err := this.dial_getinfo()
	util.Try(err)
	rst, err := this.dial_login(info)
	util.Try(err)
	util.DEBUG.Log(rst)
	return
}

type loginInfo struct {
	UserIP   string `xml:"Redirect>UserIP"`
	LoginURL string `xml:"Redirect>LoginURL"`
	Uuid     string `xml:"Redirect>Uuid"`
}
type loginResult struct {
	ResponseCode string `xml:"AuthenticationReply>ResponseCode"`
	LogoffURL    string `xml:"AuthenticationReply>LogoffURL"`
	Uuid         string `xml:"AuthenticationReply>Uuid"`
	UserIP       string `xml:"AuthenticationReply>UserIP"`
}

func (this *Dialer) dial_getinfo() (info *loginInfo, err error) {
	defer util.Catch(&err)
	body := ioutil.NopCloser(strings.NewReader((url.Values{"wlanuserip": {this.UserIP}}).Encode()))
	req, err := http.NewRequest("POST", "http://115.239.134.163:8080/showlogin.do", body)
	util.Try(err)
	req.Header.Set("User-Agent", "China Telecom Client")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	util.Try(err)
	de := xml.NewDecoder(res.Body)
	info = &loginInfo{}
	de.Decode(info)
	res.Body.Close()
	return
}
func (this *Dialer) dial_login(info *loginInfo) (rst *loginResult, err error) {
	defer util.Catch(&err)
	body := ioutil.NopCloser(strings.NewReader((url.Values{
		"uuid":       {info.Uuid},
		"userip":     {info.UserIP},
		"username":   {this.username},
		"password":   {this.rawPassword},
		"ratingtype": {this.ratingtype},
	}).Encode()))
	req, err := http.NewRequest("POST", info.LoginURL, body)
	util.Try(err)
	req.Header.Set("User-Agent", "China Telecom Client")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	util.Try(err)
	de := xml.NewDecoder(res.Body)
	rst = &loginResult{}
	de.Decode(rst)
	res.Body.Close()
	return
}

func (this *Dialer) SetPassword(pwd string) {
	this.rawPassword = pwd
	this.password = util.AESCBCStringX(util.Md5(NET_PWD_MASK), []byte(pwd), true)
}

func (this *Dialer) RefreshIP() {
	for _, v := range util.GetIPLocal() {
		if v[:3] == "10." {
			this.UserIP = v
			return
		}
	}
}

// 路由拨号加密用户名
func (this *Dialer) getCryptUsername() string {
	time := util.JsCurrentSecond() / 5
	data := string([]rune{
		rune(time >> 24 & 0xff),
		rune(time >> 16 & 0xff),
		rune(time >> 8 & 0xff),
		rune(time & 0xff),
	})
	data += this.username[:strings.IndexRune(this.username, '@')]
	data += RADIUS
	aftermd5 := util.Md5String(data)
	util.DEBUG.Log(aftermd5)
	sig := aftermd5[:2]
	temp := make([]byte, 32)
	timechar := []byte{
		byte(time >> 24 & 0xff),
		byte(time >> 16 & 0xff),
		byte(time >> 8 & 0xff),
		byte(time & 0xff),
	}
	for i := 0; i < 32; i++ {
		temp[i] = timechar[(31-i)>>3] & 1
		timechar[(31-i)>>3] = timechar[(31-i)>>3] >> 1
	}
	timeHash := make([]byte, 4)
	for i := 0; i < 4; i++ {
		timeHash[i] = temp[i]<<7 + temp[4+i]<<6 + temp[8+i]<<5 + temp[12+i]<<4 + temp[16+i]<<3 + temp[20+i]<<2 + temp[24+i]<<1 + temp[28+i]
	}
	temp[0] = (timeHash[0] >> 2) & 0x3F
	temp[1] = (timeHash[0]&3)<<4 + (timeHash[1] >> 4 & 0xF)
	temp[2] = ((timeHash[2] >> 6) & 0x3) + (timeHash[1]&0xF)<<2
	temp[3] = timeHash[2] & 0x3F
	temp[4] = (timeHash[3] >> 2) & 0x3F
	temp[5] = (timeHash[3] & 3) << 4
	sig2 := ""
	for i := 0; i < 6; i++ {
		var tp = temp[i] + 0x20
		if tp >= 0x40 {
			tp++
		}
		sig2 += string([]byte{byte(tp)})
	}
	return "\r\n" + sig2 + sig + this.username
}
