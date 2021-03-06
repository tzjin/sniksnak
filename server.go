package main

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/context"

	"sniksnak/controllers"
	"sniksnak/system"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func main() {

	defer glog.Flush()

	var application = &system.Application{}

	application.Init()
	application.LoadTemplates()

	// Setup static files
	static := web.New()
	publicPath := "public"
	static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(publicPath))))

	http.Handle("/assets/", static)

	// Apply middleware
	goji.Use(application.ApplyTemplates)
	goji.Use(application.ApplyDbMap)
	goji.Use(context.ClearHandler)

	controller := &controllers.MainController{}
	apicontroller := &controllers.ApiController{}

	// Couple of files - in the real world you would use nginx to serve them.
	goji.Get("/assets/robots.txt", http.FileServer(http.Dir(publicPath)))
	goji.Get("/assets/img/favicon.ico", http.FileServer(http.Dir(publicPath+"/images")))

	// Home page
	goji.Get("/", application.Route(controller, "Index"))

	// handlers for /api/* calls
	goji.Get("/api/get/", apicontroller.GET_data)
	goji.Post("/api/inc/:id", apicontroller.INC_counter)
	goji.Post("/api/dec/:id", apicontroller.DEC_counter)

	graceful.PostHook(func() {
		application.Close()
	})

	goji.Serve()
}
