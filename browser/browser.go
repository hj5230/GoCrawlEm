package browser

import (
	"context"

	"github.com/chromedp/chromedp"
)

func CreateContext(iCtx context.Context, opts ...chromedp.ContextOption) (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewContext(iCtx, opts...)
	return ctx, cancel
}

func AllocateContext(url string) (context.Context, context.CancelFunc) {
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), url)
	return allocatorCtx, cancel
}
