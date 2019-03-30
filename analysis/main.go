package main

import (
	"flag"
	"fmt"
	"time"
)

type cmdParams struct {
	logFilePath string
	routineNum int
}

type digData struct {
	time string
	url string
	refer string
	ua string
}

type urlData struct {
	data digData
	uid string
}

type urlNode struct {

}

type storageBlock struct {
	counterType string
	storageModel string
	unode urlNode
}

func init()  {

}

// 获取参数->打日志->初始化channel->日志消费->创建一组日志处理->创建pv uv->创建存储器
func main() {

	logFilePath := flag.String("filePath", "/Users/web/go/log/miss-log.log", "log file path ")
	routineNum := flag.Int("routineNum", 5, "consumer number by goroutine")
	//l := flag.String("l", "/Users/web/go/log", "this is go log file path")
	flag.Parse()

	params := cmdParams{*logFilePath, *routineNum}
	fmt.Println(params)


	logChannel := make(chan string, 3*params.routineNum)
	pvChannel := make(chan urlData, params.routineNum)
	uvChannel := make(chan urlData, params.routineNum)
	storageChannel := make(chan storageBlock, params.routineNum)

	// 日志消费者
	go readFileLineByLine(params, logChannel)

	// 创建日志处理
	for i := 0; i <params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	// 创建pv uv统计
	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel)

	// 创建存储器
	go dataStorage(storageChannel)

	fmt.Printf("ceshi")

	time.Sleep(1000*time.Second)
}

func dataStorage(storageChannel chan storageBlock)  {
	
}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock)  {

}

func uvCounter(pvChannel chan urlData, storageChannel chan storageBlock)  {

}

func logConsumer(logChannel chan string, pvChannel chan urlData, uvChannel chan urlData)  {

}

func readFileLineByLine(params cmdParams, logChannel chan string)  {
	
}
