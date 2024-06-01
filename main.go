package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hj5230/GoCrawlEm/browser"
	"github.com/hj5230/GoCrawlEm/gocrawlem"
)

const (
	url        = "https://peoplesdaily.pdnews.cn/searchDetails"
	searchTerm = "Huawei"
)

func main() {
	gocrawlem.Gocrawlem()

	alloc, cancel := browser.AllocateContext("ws://localhost:9222") // use docker context
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
		// Get all the targets (open tabs)
		targets, err := chromedp.Targets(ctx)
		if err != nil {
			return err
		}

		// Find the new tab (the one that isn't the current tab)
		for _, target := range targets {
			if target.Type == "page" && target.URL != url {
				newTabCtx, cancel = browser.CreateContext(ctx, chromedp.WithTargetID(target.TargetID))
				return nil // Found the target, exit the loop
			}
		}

		return nil
	}))
	if err != nil {
		panic(err)
	}
	defer cancel() // Defer the cancel function for the new tab context

	// Switch to the new tab and extract the news title
	err = chromedp.Run(newTabCtx,
		// Wait for the news title to be visible
		chromedp.WaitVisible("#newsTitle", chromedp.ByID),

		// Extract the content of the newsTitle div
		chromedp.Text("#newsTitle", &newsTitle, chromedp.ByID),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("First News Title:", newsTitle)
}
