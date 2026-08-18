package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
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

// mongoURL comes from the MONGO_URI environment variable
var mongoURL = getEnv("MONGO_URI", "mongodb://localhost:27017")

// getEnv reads an environment variable, or returns a fallback if it's unset.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type Employee struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name"`
	Salary float64            `json:"salary"`
	Age    int                `json:"age"`
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

	// Logs every request (method, path, status, latency) to stdout.
	app.Use(logger.New())

	app.Get("/employee", func(c fiber.Ctx) error {
		// Give every DB call a timeout so a stuck DB can't hang the request forever.
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		cursor, err := mg.DB.Collection("employees").Find(ctx, bson.D{})
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var employee []Employee = make([]Employee, 0)

		if err := cursor.All(ctx, &employee); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(employee)

	})

	app.Get("/employee/:id", func(c fiber.Ctx) error {
		employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		query := bson.D{{Key: "_id", Value: employeeID}}

		employee := new(Employee)
		err = mg.DB.Collection("employees").FindOne(ctx, query).Decode(employee)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).SendString("employee not found")
			}
			return c.Status(500).SendString(err.Error())
		}

		return c.Status(200).JSON(employee)
	})

	app.Post("/employee", func(c fiber.Ctx) error {
		collection := mg.DB.Collection("employees")

		employee := new(Employee)

		if err := c.Bind().Body(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		// Basic sanity checks so obviously-bad data can't be inserted.
		if employee.Name == "" {
			return c.Status(400).SendString("name is required")
		}
		if employee.Salary < 0 {
			return c.Status(400).SendString("salary cannot be negative")
		}
		if employee.Age <= 0 {
			return c.Status(400).SendString("age must be greater than 0")
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		insertResult, err := collection.InsertOne(ctx, employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		employee.ID = insertResult.InsertedID.(primitive.ObjectID)

		return c.Status(201).JSON(employee)
	})

	app.Put("/employee/:id", func(c fiber.Ctx) error {
		employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		employee := new(Employee)
		if err := c.Bind().Body(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		if employee.Name == "" {
			return c.Status(400).SendString("name is required")
		}
		if employee.Salary < 0 {
			return c.Status(400).SendString("salary cannot be negative")
		}
		if employee.Age <= 0 {
			return c.Status(400).SendString("age must be greater than 0")
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		query := bson.D{{Key: "_id", Value: employeeID}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "salary", Value: employee.Salary},
				{Key: "age", Value: employee.Age},
			}},
		}

		err = mg.DB.Collection("employees").FindOneAndUpdate(ctx, query, update).Err()
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

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		query := bson.D{{Key: "_id", Value: employeeID}}
		result, err := mg.DB.Collection("employees").DeleteOne(ctx, &query)
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
