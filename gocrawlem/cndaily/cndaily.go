package cndaily

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/chromedp/chromedp"
	"github.com/hj5230/GoCrawlEm/browser"
)

const (
	url   = "https://newssearch.chinadaily.com.cn/en/search?cond=%7B%22publishedDateFrom%22%3A%222020-06-30%22%2C%22publishedDateTo%22%3A%222024-06-01%22%2C%22titleMust%22%3A%22Huawei%22%2C%22channel%22%3A%5B%222%40cndy%22%2C%222%40webnews%22%2C%222%40bw%22%2C%222%40hk%22%2C%22ismp%40cndyglobal%22%5D%2C%22type%22%3A%5B%22story%22%2C%22comment%22%2C%22blog%22%5D%2C%22curType%22%3A%22story%22%2C%22sort%22%3A%22dp%22%2C%22duplication%22%3A%22on%22%7D&language=en"
	fPath = "./gocrawlem/cndaily/urls.json"
)

type PageURLs struct {
	PageUrls [][]string `json:"pageUrls"`
}

func CrawlPostUrls() []string {
	dockerCtx, cancelDocker := browser.UseDockerContext() // Use Docker context
	defer cancelDocker()

	ctx, cancel := browser.CreateContext(dockerCtx) // Create a new browser context
	defer cancel()

	page := 1
	var allHrefs []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Printf("Crawling page: %d\n", page)

		var hrefs []string
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`Array.from(document.querySelectorAll('.intro a')).map(a => a.href);`, &hrefs),
		)
		if err != nil {
			panic(err)
		}
		allHrefs = append(allHrefs, hrefs...)

		fmt.Printf("Found %d urls on page %d\n", len(hrefs), page)

		var nextExists bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelector('.page.rt a[title="next"]') !== null`, &nextExists),
		)
		if err != nil {
			panic(err)
		}
		if !nextExists {
			break
		} else {
			fmt.Printf("Next page is %d\n", page+1)
		}

		err = chromedp.Run(ctx,
			chromedp.Click(`a[title="next"]`, chromedp.NodeVisible, chromedp.ByQuery),
			chromedp.Sleep(3*time.Second),
		)
		if err != nil {
			panic(err)
		}
		page++
		fmt.Printf("Arrived at page %d successfully\n", page)
	}

	return allHrefs
}

func CrawlPostContent(url string) (string, string, string) {
	dockerCtx, cancelDocker := browser.UseDockerContext() // Use Docker context
	defer cancelDocker()

	ctx, cancel := browser.CreateContext(dockerCtx) // Create a new browser context
	defer cancel()

	var title string
	var info string
	var content string

	defer func() {
		if r := recover(); r != nil {
			title, info, content = url, "", ""
		}
	}()

	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`
			document.querySelector('h1').innerText;
		`, &title),
		chromedp.Evaluate(`
			document.querySelector('.info_l').innerText;
		`, &info),
		chromedp.Evaluate(`
			(() => {
				let content = '';
				document.querySelectorAll('#Content p').forEach(p => {
					content += p.innerText + ' ';
				});
				return content;
			})();
		`, &content),
	)
	if err != nil {
		return url, "", ""
	}

	return title, info, content
}

func CrawlFromJson() {
	fData, err := os.ReadFile(fPath)
	if err != nil {
		panic(err)
	}

	var urls PageURLs

	err = sonic.Unmarshal(fData, &urls)
	if err != nil {
		panic(err)
	}

	allUrls := make([]string, 0)
	for _, pageUrls := range urls.PageUrls {
		allUrls = append(allUrls, pageUrls...)
	}

	fmt.Println(len(allUrls), "urls have been loaded")

	fi, err := os.Create("cndaily.csv")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	writer := csv.NewWriter(fi)
	defer writer.Flush()

	header := []string{"Title", "Info", "Content"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	fmt.Println("Header attached")

	for i, url := range allUrls {
		fmt.Println("Crawling url", i+1, "out of", len(allUrls))

		title, info, content := CrawlPostContent(url)

		content = strings.ReplaceAll(content, "\n", " ")

		row := []string{title, info, content}
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}
}
