package main

import (
	"encoding/json"
	"fmt"
	"go_log/collect_log/dbops"
	"go_log/collect_log/defs"
	"go_log/collect_log/util"
	"time"
)

func callback(d util.MSG) {
	fmt.Println("Ok")
	var ubody defs.Log
	data := d.Body
	json.Unmarshal([]byte(data), &ubody)
	fmt.Println(ubody)
	if err := dbops.AddLog(&ubody); err != nil {
		return
	}


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
