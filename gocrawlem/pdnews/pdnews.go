package pdnews

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hj5230/GoCrawlEm/browser"
)

const (
	url        = "https://peoplesdaily.pdnews.cn/searchDetails"
	searchTerm = "Huawei"
)

func Crawl() {
	alloc, cancel := browser.AllocateDockerContext() // use docker context
	defer cancel()
	ctx, cancel := browser.CreateContext(alloc) // create a new context
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(".search-input", chromedp.ByQuery),                // find the search input
		chromedp.SendKeys(".search-input input", searchTerm, chromedp.ByQuery), // enter the search term
		chromedp.Click(".search-input_btn", chromedp.ByQuery),                  // click the search button
		chromedp.Click(".search-details_item:first-child", chromedp.ByQuery),   // click the first result for test
	)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	var newsTitle string
	var newTabCtx context.Context

	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		targets, err := chromedp.Targets(ctx)
		if err != nil {
			return err
		}

		for _, target := range targets {
			if target.Type == "page" && target.URL != url {
				newTabCtx, cancel = browser.CreateContext(ctx, chromedp.WithTargetID(target.TargetID))
				return nil
			}
		}

		return nil
	}))
	if err != nil {
		panic(err)
	}
	defer cancel()

	err = chromedp.Run(newTabCtx,
		chromedp.WaitVisible("#newsTitle", chromedp.ByID),      // wait for the news title to be visible
		chromedp.Text("#newsTitle", &newsTitle, chromedp.ByID), // get the news title
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("First News Title:", newsTitle)
}
