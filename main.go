package main

import (
	"log"
	"os"

	"price-tracker/parser"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Server is worked!"})
	})

	app.Get("/parse", func(c *fiber.Ctx) error {
		url := "https://indexiq.ru/catalog/iphone-15-pro-max/"
	
		products, err := parser.ScrapeProducts(url)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	
		return c.JSON(products)
	})
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
