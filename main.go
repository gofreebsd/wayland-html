package main

import "C"

import (
	"code.google.com/p/go.net/websocket"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func http() {
	server := martini.Classic()

	server.Use(martini.Static("public", martini.StaticOptions{
		Prefix: "public",
	}))

	server.Use(martini.Static("bower_components", martini.StaticOptions{
		Prefix: "bower_components",
	}))

	server.Use(render.Renderer(render.Options{
		Extensions: []string{".tmpl", ".html"},
		Delims: render.Delims{"{[{", "}]}"},
	}))

	server.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	//
	server.Get("/clients/:pid", func(params martini.Params, r render.Render) {
		r.HTML(200, "index", params["pid"])
	})

	type jd map[string]interface{}
	type ja []interface{}

	// json api
	server.Get("/api/clients", func(r render.Render) {
		client_info := ja{}

		for _, c := range compositors {
			client_info = append(client_info, jd{
				"pid": c.Pid,
			})
		}

		r.JSON(200, client_info)
	})

	// websocket api

	server.Get("/api/clients", websocket.Handler(func(ws *websocket.Conn) {

	}).ServeHTTP)

	server.Get("/api/clients/:pid", websocket.Handler(func(ws *websocket.Conn) {

	}).ServeHTTP)

	server.Run()
}

func main() {

	go http()

	wayland()
}
