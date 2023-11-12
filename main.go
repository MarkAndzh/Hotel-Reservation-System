package main

import (
	"context"
	"flag"
	"github.com/MarkAndzh/hotel-reservation/api"
	"github.com/MarkAndzh/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"os"
)

const dburi = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	listenAddr := flag.String("listenAddr", ":5001", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		logger.Error(err.Error())
	}
	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(logger, client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	err = app.Listen(*listenAddr)
	if err != nil {
		panic(err)
	}
}
