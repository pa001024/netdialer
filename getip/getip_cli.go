package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pa001024/reflex/util"
)

func main() {
	check()
}

func check() {
	ch := make(chan string)
	go func() {
		ip := util.GetIPInfo()
		if ip != nil {
			ch <- ip.IP
		} else {
			os.Exit(2)
		}
	}()
	select {
	case <-time.After(time.Second * 2):
		os.Exit(2)
	case ip := <-ch:
		fmt.Println(ip)
	}
}
