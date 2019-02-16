package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jonnylangefeld/go-api-base-project/model"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var db *gorm.DB

func TestMain(m *testing.M)  {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it with the given run parameters
	opts := dockertest.RunOptions{
		Repository:   "mdillon/postgis",
		Tag:          "latest",
		Env:          []string{"POSTGRES_PASSWORD=mysecretpassword"},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {{HostIP: "0.0.0.0", HostPort: "5433"}},
		},
	}
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Exponential retry to connect to database while it is booting
	if err := pool.Retry(func() error {
		databaseConnStr := fmt.Sprintf("host=localhost port=5433 user=postgres dbname=postgres password=mysecretpassword sslmode=disable")
		_ = os.Setenv("DATABASE_SERVER", databaseConnStr)
		db, err = gorm.Open("postgres", databaseConnStr)

		// if there was an issue opening the connection
		if err != nil {
			log.Println("Database not ready yet (it is booting up, wait for a few tries)...")
			return err
		}

		return db.DB().Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	log.Println("Initialize test database...")
	initTestDatabase()

	// Run the actual test cases (functions that start with Test...)
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestName(t *testing.T)  {
	app := newApp(db)
	e := httptest.New(t, app, httptest.URL("http://localhost"))

	t1 := e.GET("/Bill").Expect().Status(iris.StatusOK)
	assert.Equal(t, "Hello Bill", t1.Body().Raw())
}

func TestOrders(t *testing.T)  {
	app := newApp(db)
	e := httptest.New(t, app, httptest.URL("http://localhost"))

	t1 := e.GET("/orders").Expect().Status(iris.StatusOK)
	expected, _ := json.Marshal(sampleOrders)
	assert.Equal(t, string(expected), t1.Body().Raw())
}

func initTestDatabase()  {
	db.AutoMigrate(&model.Order{})

	db.Save(&sampleOrders[0])
	db.Save(&sampleOrders[1])
}

var sampleOrders = []model.Order{
		{
			ID: 1,
			Description: "An old glove",
			Ts: time.Now().Unix() * 1000,
		},
		{
			ID: 2,
			Description: "Something you don't need",
			Ts: time.Now().Unix() * 1000,
		},
	}