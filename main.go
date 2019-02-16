package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jonnylangefeld/go-api-base-project/model"
	"github.com/jonnylangefeld/go-api-base-project/util"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main()  {

	// Create the databse connection
	db, err := gorm.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=postgres password=mysecretpassword sslmode=disable",
		)
	util.RaiseFatalErrorIf(err, "could not establish connection to database.")

	db.AutoMigrate(&model.Order{})

	// Closes the database connection when the program ends
	defer func() {
		_ = db.Close()
	}()

	// Initialize the app
	app := newApp(db)

	err = app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	util.RaiseFatalErrorIf(err, "Could not start iris app.")
}

func newApp(db *gorm.DB) *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/", func(ctx iris.Context){
		_,_ = ctx.Text("Hello World")
	})

	app.Get("{name:string}", func(ctx iris.Context) {
		_,_ = ctx.Text(fmt.Sprintf("Hello %s", ctx.Params().Get("name")))
	})

	var orders []model.Order

	app.Get("/orders", func(ctx iris.Context) {
		db.Find(&orders)
		_,_ = ctx.JSON(orders)
	})

	return app
}