package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/pa001024/netdialer"
	"github.com/pa001024/reflex/util"
)

func main() {
	username := flag.String("u", "", "shanxun username")
	password := flag.String("p", "", "shanxun password")
	isRouter := flag.Bool("r", false, "use router?")
	realIP := flag.String("ip", "local", "real IP, 'stdin' or 'local' works")
	flag.Parse()
	if len(*username) < 12 || len(*password) != 6 {
		util.ERROR.Log("Invalid username or password")
		flag.Usage()
		return
	}
	d := netdialer.NewDialer(*username, *password)
	if *realIP != "" && *realIP != "local" {
		if *realIP == "stdin" {
			in := bufio.NewReader(os.Stdin)
			str, _ := in.ReadString(0)
			d.UserIP = str
		} else {
			d.UserIP = *realIP
		}
	}
	var err error
	if *isRouter {
		err = d.ConnectRouter()
	} else {
		err = d.ConnectDirect()
	}
	if err != nil {
		util.ERROR.Log(err)
	}
}
