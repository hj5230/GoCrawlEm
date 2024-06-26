package browser

import (
	"context"

	"github.com/chromedp/chromedp"
)

const (
	docker = "ws://localhost:9222"
)

func CreateContext(iCtx context.Context, opts ...chromedp.ContextOption) (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewContext(iCtx, opts...)
	return ctx, cancel
}

// deprecated
func AllocateDockerContext() (context.Context, context.CancelFunc) {
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), docker)
	return allocatorCtx, cancel
}

func UseDockerContext() (context.Context, context.CancelFunc) {
	dockerCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), docker)
	return dockerCtx, cancel
}
