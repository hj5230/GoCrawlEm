package main

import (
	"github.com/hj5230/GoCrawlEm/gocrawlem/cndaily"
)

func main() {
	// xinhuar.Crawl()
	// title, info, content := cndaily.CrawlPostContent("https://www.chinadaily.com.cn/a/202403/07/WS65e928cba31082fc043bb29e.html")
	// fmt.Println(title, info)
	// fmt.Println(content)
	cndaily.CrawlFromJson()
}
