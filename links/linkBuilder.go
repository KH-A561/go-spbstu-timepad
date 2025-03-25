package links

import (
	"context"
	"github.com/valyala/fasthttp"
)

type LinkBuilder[E any] interface {
	Fetch(ctx context.Context) (E, error)
}

type DefaultFetcherImpl[E any] struct {
	Client *fasthttp.Client
}
