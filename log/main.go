package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
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

func main() {
	total := flag.Int("total", 100, "how many rows by created")
	filePath := flag.String("filePath", "/Users/web/go/log", "log file path ")
	flag.Parse()

	res := ruleResource()
	list := buildUrl(res)

	fmt.Println(*total, *filePath)

	fmt.Println(list)

	//for i := 0; i <= *total; i ++ {
	//	logStr := ""
	//	ioutil.WriteFile(*filePath, []byte(logStr), 0777)
	//}

}