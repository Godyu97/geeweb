package main

import (
	"fmt"
	"github.com/Godyu97/geeweb/gee"
)

func main() {
	mux := gee.New()
	mux.GET("/", func(ctx *gee.Context) {
		fmt.Fprintf(ctx.Writer, "URL.Path = %q\n", ctx.Path)
	})
	mux.GET("/:hello", func(ctx *gee.Context) {
		for k, v := range ctx.Req.Header {
			fmt.Fprintf(ctx.Writer, "Header[%q] = %q\n", k, v)
		}
	})
	mux.Run(":9999")
}
