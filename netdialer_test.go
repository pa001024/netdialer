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
	dialer = NewDialer("17805805321@HYXY.XY", "123456")
}

func TestGetCryptUsername(t *t.T) {
	out := dialer.getCryptUsername()
	fmt.Println(out)
}
