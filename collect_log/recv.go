package main

import (
	"fmt"
	"go_log/collect_log/defs"
	"go_log/collect_log/util"
	"time"
)

func callback(d util.MSG) {
	fmt.Println("Ok")
	fmt.Println(string(d.Body))
}

func errCallback(d util.MSG) {
	fmt.Println("Err")
	fmt.Println(string(d.Body))
}

func dlxCallback(d util.MSG) {
	fmt.Println("Dlx")
	fmt.Println(string(d.Body))
}

func main() {
	if err := util.Init(defs.CONFIG_PATH + "./rmq.json"); err != nil {
		fmt.Println(err)
	}

	if err := util.Pop("myPoper", callback); err != nil {
		fmt.Println(err)
	}


	if err := util.Pop("errPoper", errCallback); err != nil {
		fmt.Println(err)
	}

	if err := util.Pop("dlxPoper", dlxCallback); err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Duration(1000)*time.Second)

	if err := util.Fini(); err != nil {
		fmt.Println(err)
	}
}
