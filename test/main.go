package main

import (
	"fmt"
	"github.com/Godyu97/geeweb/gee"

	"html/template"
	"net/http"

	"time"
)

func main() {
	mux := gee.New()
	mux.GET("/", func(ctx *gee.Context) {
		fmt.Fprintf(ctx.Writer, "URL.Path = %q\n", ctx.Path)
	})
	mux.GET("/dy/:hello", func(ctx *gee.Context) {
		for k, v := range ctx.Req.Header {
			fmt.Fprintf(ctx.Writer, "Header[%q] = %q\n", k, v)
		}
	})
	//static fileserver
	g := mux.Group("/v1")
	g.Static("/assets", "./static/assets")
	//html template
	mux.SetHtmlFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	mux.LoadHTMLGlob("./static/templates/*")
	stu1 := student{Name: "Geektutu", Age: 20}
	stu2 := student{Name: "Jack", Age: 22}
	mux.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]student{stu1, stu2},
		})
	})

	mux.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Now(),
		})
	})
	mux.GET("/css", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	mux.Run(":9999")
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	return t.String()
}
