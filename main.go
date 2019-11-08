package main

import (
	"fmt"
	"my-go-api/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kataras/iris/v12"
)

func main() {

	// Create the databse connection
	db, err := gorm.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=postgres password=mysecretpassword sslmode=disable",
	)
	// End the program with an error if it could not connect to the database
	if err != nil {
		panic("could not connect to database")
	}

	// Create the default database schema and a table for the orders
	db.AutoMigrate(&model.Order{})

	// Closes the database connection when the program ends
	defer func() {
		_ = db.Close()
	}()

	// Run the function to create the new iris App
	app := newApp(db)

	// Start the web server on port 8080
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func newApp(db *gorm.DB) *iris.Application {
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

	// Define the slice for the result
	var orders []model.Order

	// Endpoint to perform the database request
	app.Get("/orders", func(ctx iris.Context) {
		db.Find(&orders)
		_, _ = ctx.JSON(orders)
	})

	return app
}
