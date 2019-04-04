package taskrunner

import (
	"encoding/json"
	"fmt"
	"go_log/collect_log/dbops"
	"go_log/collect_log/defs"
	"go_log/collect_log/util"
)

func Start()  {
	if err := util.Pop("myPoper", callback1); err != nil {
		fmt.Println(err)
	}


	if err := util.Pop("errPoper", errCallback1); err != nil {
		fmt.Println(err)
	}

	if err := util.Pop("dlxPoper", dlxCallback1); err != nil {
		fmt.Println(err)
	}
}

func callback1(d util.MSG) {
	//fmt.Println("Ok")
	var ubody defs.Log
	data := d.Body
	json.Unmarshal([]byte(data), &ubody)
	//fmt.Println(ubody)
	if err := dbops.AddLog(&ubody); err != nil {
		return
	}


}

func errCallback1(d util.MSG) {
	fmt.Println("Err")
	fmt.Println(string(d.Body))
}

func dlxCallback1(d util.MSG) {
	//fmt.Println("Dlx")
	var ubody defs.Log
	//data := d.Body
	json.Unmarshal([]byte(data), &ubody)
	fmt.Println(ubody)
	if err := dbops.AddLog(&ubody); err != nil {
		return
	}
}
