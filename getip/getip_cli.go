package main

import (
	"fmt"
	"os"

	"github.com/pa001024/reflex/util"
)

func main() {
	check()
}

func check() {
	ip := util.GetIPInfo()
	if ip != nil {
		fmt.Println(ip.IP)
	} else {
		os.Exit(2)
	}
}
