package cndaily

import (
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hj5230/GoCrawlEm/browser"
)

const (
	url = "https://newssearch.chinadaily.com.cn/en/search?cond=%7B%22publishedDateFrom%22%3A%222020-06-30%22%2C%22publishedDateTo%22%3A%222024-06-01%22%2C%22titleMust%22%3A%22Huawei%22%2C%22channel%22%3A%5B%222%40cndy%22%2C%222%40webnews%22%2C%222%40bw%22%2C%222%40hk%22%2C%22ismp%40cndyglobal%22%5D%2C%22type%22%3A%5B%22story%22%2C%22comment%22%2C%22blog%22%5D%2C%22curType%22%3A%22story%22%2C%22sort%22%3A%22dp%22%2C%22duplication%22%3A%22on%22%7D&language=en"
)

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
		log.Fatalf("Failed to navigate to the URL: %v", err)
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
