package cndaily

import (
	"github.com/chromedp/chromedp"
	"github.com/hj5230/GoCrawlEm/browser"
)

const (
	url = "https://newssearch.chinadaily.com.cn/en/search?cond=%7B%22publishedDateFrom%22%3A%222020-06-30%22%2C%22publishedDateTo%22%3A%222024-06-01%22%2C%22titleMust%22%3A%22Huawei%22%2C%22channel%22%3A%5B%222%40cndy%22%2C%222%40webnews%22%2C%222%40bw%22%2C%222%40hk%22%2C%22ismp%40cndyglobal%22%5D%2C%22type%22%3A%5B%22story%22%2C%22comment%22%2C%22blog%22%5D%2C%22curType%22%3A%22story%22%2C%22sort%22%3A%22dp%22%2C%22duplication%22%3A%22on%22%7D&language=en"
)

func CrawlPostUrls() []string {
	alloc, cancel := browser.AllocateDockerContext() // use docker context
	defer cancel()
	ctx, cancel := browser.CreateContext(alloc) // create a new context
	defer cancel()

	var hrefs []string
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('h4 a')).map(a => a.href);
		`, &hrefs), // then next page... repeat until the end
	)
	if err != nil {
		panic(err)
	}

	return hrefs
}
