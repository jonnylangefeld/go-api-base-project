package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kataras/iris"
)

func main() {

	// Run the function to create the new iris App
	app := newApp()

	// Start the web server on port 8080
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func newApp() *iris.Application {
	// Initialize a new iris App
	app := iris.New()

	// Register the request handler for the endpoint "/"
	app.Get("/", func(ctx iris.Context) {
		// Return something by adding it to the context
		ctx.Text("Hello World")
	})

	// Register an endpoint with a variable
	app.Get("{name:string}", func(ctx iris.Context) {
		_, _ = ctx.Text(fmt.Sprintf("Hello %s", ctx.Params().Get("name")))
	})

	return app
}
