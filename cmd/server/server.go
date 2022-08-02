package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

// request handler in fasthttp style, i.e. just plain function.
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
}

func main() {
	fasthttp.ListenAndServe(":8087", fastHTTPHandler)
}
