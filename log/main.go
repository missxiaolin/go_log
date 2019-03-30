package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type resource struct {
	url string
	target string
	start int
	end int
}

func ruleResource() []resource {
	var res []resource
	r1 := resource{
		url: "http://www.miss-log.com",
		target: "",
		start: 0,
		end: 0,
	}
	r2 := resource{
		url: "http://www.miss-log.com/list/{$id}.html",
		target: "{$id}",
		start: 1,
		end: 21,
	}
	r3 := resource{
		url: "http://www.miss-log.com/detail/{$id}.html",
		target: "{$id}",
		start: 1,
		end: 10924,
	}
	res = append(append(append(res, r1), r2), r3)
	return res
}

func buildUrl(res []resource) []string {
	var list [] string

	for _,v := range res {
		if len(v.target) == 0 {
			list = append(list, v.url)
		} else {
			for i := v.start; i <= v.end; i++ {
				urlString := strings.Replace(v.url, v.target, strconv.Itoa(i), -1)
				list = append(list, urlString)
			}
		}

	}
	return list
}

func makLog(current, refer, ua string) string {
	u := url.Values{}
	u.Set("time", "1")
	u.Set("url", current)
	u.Set("refer", refer)
	u.Set("ua", ua)
	paramStr := u.Encode()
	logTem := "127.0.0.1 - - [28/Mar/2019:10:53:40 +0800] \"GET /dig?{$paramStr} HTTP/1.1\" 200 43 \"-\" \"{$ua}\""
	log := strings.Replace(logTem, "{$paramStr}" ,paramStr, -1)
	log = strings.Replace(log, "{$ua}", ua, -1)
	return log
}

func randInt(min, max int) int {
	// 取随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if min > max {
		return max
	}
	return r.Intn(max-min) + min
}

func main() {
	total := flag.Int("total", 100, "how many rows by created")
	filePath := flag.String("filePath", "/Users/web/go/log/miss-log.log", "log file path ")
	flag.Parse()

	res := ruleResource()
	list := buildUrl(res)

	//fmt.Println(*total, *filePath)

	//fmt.Println(list)

	logStr := ""

	for i := 0; i <= *total; i ++ {
		currentUrl := list[randInt(0, len(list)-1)]
		referUrl := list[randInt(0, len(list)-1)]
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36"
		logStr =logStr + makLog(currentUrl, referUrl, ua) + "\n"
		//ioutil.WriteFile(*filePath, []byte(logStr), 0777)
		//fd,_ := os.OpenFile(*filePath, os.O_RDWR|os.O_APPEND,0644)
	}
	fd,_ := os.OpenFile(*filePath, os.O_RDWR|os.O_APPEND,0644)
	fd.Write([]byte(logStr))
	defer fd.Close()

	fmt.Printf("done.\n")
}