package main

import (
	//"fmt"
	"github.com/codegangsta/martini"
	//"github.com/elvtechnology/gocqltable"
	//"github.com/martini-contrib/render"
	//"net/http"
)

func main() {
	m := martini.Classic()
	m.Get("/hello", func() string {
		return "hello world"
	})
	m.Run()
}
