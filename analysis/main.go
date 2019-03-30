package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mgutz/str"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const HANDLE_DIG  = " /dig?"
const HANDLE_MOVIE = "/movie/"
const HANDLE_LIST = "/list/"
const HANDLE_HTML  = ".html"

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
	unode  urlNode
}

type urlNode struct {
	unType string
	unRid int
	unUrl string
	unTime string
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

	// Redis Pool
	redisPool, err := pool.New( "tcp", "127.0.0.1:6379", 2*params.routineNum )
	if err != nil{
		fmt.Printf( "Redis pool created failed." )
		panic(err)
	} else {
		go func() {
			for {
				redisPool.Cmd("PING")
				time.Sleep(3 * time.Second)
			}
		}()
	}
	// 日志消费者
	go readFileLineByLine(params, logChannel)

	// 创建日志处理
	for i := 0; i <params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	// 创建pv uv统计
	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel, redisPool)

	// 创建存储器
	go dataStorage(storageChannel)

	fmt.Printf("ceshi")

	time.Sleep(1000*time.Second)
}

func dataStorage(storageChannel chan storageBlock)  {
	
}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock)  {
	for data := range pvChannel{
		sItem := storageBlock{ "pv", "ZINCRBY", data.unode }
		storageChannel <- sItem
	}
}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock, redisPool *pool.Pool)  {
	for data := range uvChannel {
		//HyperLoglog redis
		hyperLogLogKey := "uv_hpll_" + getTime(data.data.time, "day")
		ret, err := redisPool.Cmd("PFADD", hyperLogLogKey, data.uid, "EX", 86400).Int()
		if err != nil {
			fmt.Printf("UvCounter check redis hyperloglog failed, ", err)
		}
		if ret != 1 {
			continue
		}

		sItem := storageBlock{"uv", "ZINCRBY", data.unode}
		storageChannel <- sItem
	}
}

func logConsumer(logChannel chan string, pvChannel chan urlData, uvChannel chan urlData) error {
	for logStr := range logChannel {
		// 切割日志
		data := cutLogFetchData(logStr)

		hasher := md5.New()
		hasher.Write([]byte(data.refer + data.ua))
		uid := hex.EncodeToString(hasher.Sum(nil))
		uData := urlData{
			data,
			uid,
			formatUrl(data.url,data.time),
		}
		pvChannel <- uData
		uvChannel <- uData
	}

	return nil
}

func readFileLineByLine(params cmdParams, logChannel chan string) error {
	fd, err := os.Open(params.logFilePath)
	if err != nil {
		fmt.Printf("file open error")
		return err
	}
	defer fd.Close()

	count := 0
	bufferRead := bufio.NewReader(fd)
	for  {
		line, err := bufferRead.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(3*time.Second)
				fmt.Printf("no file")
			}else{
				fmt.Printf("file ReadString error")
				return err
			}

		}
		logChannel <- line
		count++

		if count%(1000*params.routineNum) == 0 {
			fmt.Printf("file read 10000")
		}
	}
	return nil
}

// 解析url
func cutLogFetchData(log string) digData {
	log = strings.TrimSpace(log)
	pos1 := str.IndexOf(log, HANDLE_DIG, 0)
	if pos1 == -1 {
		return digData{}
	}
	pos1 += len(HANDLE_DIG)

	pos2 := str.IndexOf(log, " HTTP/", pos1)
	d := str.Substr(log, pos1, pos2 - pos1)

	urlInfo, err := url.Parse("http://localhost/?" + d)
	if err != nil {
		return digData{}
	}
	data := urlInfo.Query()
	return digData{
		data.Get("time"),
		data.Get("url"),
		data.Get("refer"),
		data.Get("ua"),
	}


}

func formatUrl(url, t string) urlNode {
	pos1 := str.IndexOf( url, HANDLE_MOVIE, 0)
	if pos1!=-1 {
		pos1 += len( HANDLE_MOVIE )
		pos2 := str.IndexOf( url, HANDLE_HTML, 0 )
		idStr := str.Substr( url , pos1, pos2-pos1 )
		id, _ := strconv.Atoi( idStr )
		return urlNode{ "movie", id, url, t }
	} else {
		pos1 = str.IndexOf( url, HANDLE_LIST, 0 )
		if pos1!=-1 {
			pos1 += len( HANDLE_LIST )
			pos2 := str.IndexOf( url, HANDLE_HTML, 0 )
			idStr := str.Substr( url , pos1, pos2-pos1 )
			id, _ := strconv.Atoi( idStr )
			return urlNode{ "list", id, url, t }
		} else {
			return urlNode{ "home", 1, url, t}
		} // 如果页面url有很多种，就不断在这里扩展
	}
}


func getTime( logTime, timeType string ) string {
	var item string
	switch timeType {
	case "day":
		item = "2018-01-02"
		break
	case "hour":
		item = "2018-01-02 15"
		break
	case "min":
		item = "2018-01-02 15:04"
		break
	}
	t, _ := time.Parse(item, time.Now().Format(item))
	return strconv.FormatInt(t.Unix(), 10)

}