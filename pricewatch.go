package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))
	m.Get("/", func(ren render.Render) {
		data := make(map[string]interface{})
		data["Title"] = "Hello World"
		ren.HTML(200, "books", data)
	})
	m.Run()
}
