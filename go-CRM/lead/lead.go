package lead

import (
	"github.com/gofiber/fiber/v3"
	"github.com/soheil/go-CRM/database"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Lead struct {
	gorm.Model
	Name    string `json:"Name"`
	Company string `json:"Company"`
	Email   string `json:"Email"`
	Phone   string `json:"Phone"`
}

func GetLeads(c fiber.Ctx) error {
	var leads []Lead

	if err := database.DBConn.Find(&leads).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch leads",
		})
	}

	return c.JSON(leads)
}

func GetLead(c fiber.Ctx) error {
	id := c.Params("id")

	var lead Lead
	result := database.DBConn.First(&lead, id)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Lead not found",
		})
	}

	return c.JSON(lead)
}

func NewLead(c fiber.Ctx) error {
	lead := new(Lead)

	if err := c.Bind().Body(lead); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.DBConn.Create(lead).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(lead)
}

func DeleteLead(c fiber.Ctx) error {
	id := c.Params("id")

	var lead Lead
	result := database.DBConn.First(&lead, id)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Lead not found",
		})
	}

	if err := database.DBConn.Delete(&lead).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete lead",
		})
	}

	return c.SendStatus(204)
}
