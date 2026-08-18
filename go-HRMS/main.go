package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURL = "mongodb://localhost:27017" // or use a connection string when connecting to online services

type Employee struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name"`
	Salary float64            `json:"salary"`
	Age    float64            `json:"age"`
}

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoURL),
	)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	mg.Client = client
	mg.DB = client.Database(dbName)

	fmt.Println("Connected to MongoDB ", dbName)
	return nil
}

func main() {

	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	defer mg.Client.Disconnect(context.Background())

	app := fiber.New()

	app.Get("/employee", func(c fiber.Ctx) error {
		cursor, err := mg.DB.Collection("employees").Find(c.Context(), bson.D{})
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var employee []Employee = make([]Employee, 0)

		if err := cursor.All(c.Context(), &employee); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(employee)

	})

	// app.Get("/employee/{id}", func(c fiber.Ctx) error {

	// })

	app.Post("/employee", func(c fiber.Ctx) error {
		collection := mg.DB.Collection("employees")

		employee := new(Employee)

		if err := c.Bind().Body(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		insertResult, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		employee.ID = insertResult.InsertedID.(primitive.ObjectID)

		return c.Status(201).JSON(employee)
	})

	app.Put("/employee/{id}", func(c fiber.Ctx) error {
		employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		employee := new(Employee)
		if err := c.Bind().Body(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeID}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "salary", Value: employee.Salary},
				{Key: "age", Value: employee.Age},
			}},
		}

		err = mg.DB.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).SendString("employee not found")
			}
			return c.Status(500).SendString(err.Error())
		}

		employee.ID = employeeID

		return c.Status(200).JSON(employee)
	})

	app.Delete("/employee/:id", func(c fiber.Ctx) error {
		employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeID}}
		result, err := mg.DB.Collection("employees").DeleteOne(c.Context(), &query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if result.DeletedCount == 0 { //if no document was deleted, return 404
			return c.Status(404).SendString("employee not found")
		}

		return c.Status(200).SendString("employee deleted")
	})

	log.Fatal(app.Listen(":3000"))
}
