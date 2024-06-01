package xinhuar

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hj5230/GoCrawlEm/browser"
)

const (
	maxPage = 42
)

func CrawlPostUrls(page string) []string {
	url := fmt.Sprintf("https://so.news.cn/?lang=en#search/1/Huawei/" + page + "/0")

	fmt.Println("Crawling url page", url)

	alloc, cancel := browser.AllocateDockerContext() // use docker context
	defer cancel()
	ctx, cancel := browser.CreateContext(alloc) // create a new context
	defer cancel()

	var hrefs []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`
            Array.from(document.querySelectorAll('a')).map(a => a.href);
        `, &hrefs),
	)
	if err != nil {
		panic(err)
	}

	var filtered1 []string

	for _, href := range hrefs {
		if len(href) > 50 {
			filtered1 = append(filtered1, href)
		}
	}

	var filtered2 []string

	for i, href := range filtered1 {
		if i%2 == 0 {
			filtered2 = append(filtered2, href)
		}
	}

	return filtered2
}

func CrawlPostContent(url string) (string, string, string, string, string) {
	alloc, cancel := browser.AllocateDockerContext() // use docker context
	defer cancel()
	ctx, cancel := browser.CreateContext(alloc) // create a new context
	defer cancel()

	var title string
	var source string
	var editor string
	var time_ string
	var content string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`
			document.querySelector('h1').innerText;
		`, &title),
		chromedp.Evaluate(`
			document.querySelector('.source').innerText;
		`, &source),
		chromedp.Evaluate(`
			document.querySelector('.editor').innerText;
		`, &editor),
		chromedp.Evaluate(`
			document.querySelector('.time').innerText;
		`, &time_),
		chromedp.Evaluate(`
            (() => {
                let content = '';
                document.querySelectorAll('#detailContent p').forEach(p => {
                    let text = '';
                    p.childNodes.forEach(node => {
                        if (node.nodeType === Node.TEXT_NODE) {
                            text += node.textContent;
                        } else if (node.nodeType === Node.ELEMENT_NODE && node.tagName.toLowerCase() === 'strong') {
                            text += node.innerText;
                        }
                    });
                    if (text.trim() !== '') {
                        content += text.trim();
                    }
                });
                return content.trim();
            })();
        `, &content),
	)
	if err != nil {
		panic(err)
	}

	return title, source, editor, time_, content
}

func Crawl() {
	fi, err := os.Create("xinhuar.csv")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	writer := csv.NewWriter(fi)
	defer writer.Flush()

	header := []string{"Title", "Source", "Editor", "Time", "Content"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	for page := 0; page < maxPage; page++ {
		fmt.Println("Crawling page", page+1)

		urls := CrawlPostUrls(strconv.Itoa(page + 1))

		for _, url := range urls {
			fmt.Println("Crawling", url)

			title, source, editor, time_, content := CrawlPostContent(url)

			content = strings.ReplaceAll(content, "\n", " ")

			row := []string{title, source, editor, time_, content}

			if err := writer.Write(row); err != nil {
				panic(err)
			}
		}
	}
}
