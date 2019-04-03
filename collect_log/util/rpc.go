package util

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

var _RPCMAP map[string](map[string]func([]byte)) = make(map[string](map[string]func([]byte)))

//事件结构
type Event struct {
	FuncName string
	Params interface{}
}

//添加回调函数
func AddRpc(poper string, name string, callback func([]byte)) (err error) {
	if _, ok := _RPCMAP[poper]; !ok {
		_RPCMAP[poper] = make(map[string]func([]byte))
	}
	if _, ok := _RPCMAP[poper][name]; ok {
		return errors.New("回调函数已存在")
	} else {
		_RPCMAP[poper][name] = callback
	}
	return nil
}

//解码并回调函数
func callbackRpc(d MSG) {
	var buf bytes.Buffer
	buf.Write(d.Body)
	dec := gob.NewDecoder(&buf)
	var data Event
	dec.Decode(&data)

	if _, ok := _RPCMAP[d.Poper]; ok {
		if _, ok := _RPCMAP[d.Poper][data.FuncName]; ok {
			_RPCMAP[d.Poper][data.FuncName](data.Params.([]byte))
		}
	}
	//注，else情况还没处理

	d.Ack(false)
}


//收到Rpc调用
func RecvRpc(poperName string) (err error) {
	if err = Pop(poperName, callbackRpc); err != nil {
		return err
	}
	return nil
}


//发送RPC调用
func SendRpc(pusherName string, key string, funcName string, params interface{}) (err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	var data = Event{
		FuncName:funcName,
		Params:params,
	}
	if err = enc.Encode(&data); err != nil {
		fmt.Println("asdfasd")
		return err
	}

	if err := Push(pusherName, key, buf.Bytes()); err != nil {
		return err
	}
	return nil
}