package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/soheil/go-CRM/database"
	"github.com/soheil/go-CRM/lead"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouts(app *fiber.App) {
	app.Get("/api/v1/leads", lead.GetLeads)
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Post("/api/v1/lead", lead.NewLead)
	app.Delete("/api/v1/lead/:id", lead.DeleteLead)
}

func initDataBase() {
	var err error
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.DBConn = db

	fmt.Println("database connected")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("database migrated")
}

func main() {
	app := fiber.New() // create a new fiber app
	initDataBase()
	defer func() { // close the database connection when the program exits
		sqlDB, err := database.DBConn.DB()
		if err != nil {
			return
		}
		sqlDB.Close()
	}()
	setupRouts(app)
	app.Listen(":3000") // start the server

}
