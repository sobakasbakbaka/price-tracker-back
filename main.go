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
		shops := map[string]string{
			"indexiq": "https://indexiq.ru/catalog/iphone-15-pro-max/",
			"biggeek": "https://biggeek.ru/catalog/apple-iphone-15-pro-max",
			"store77": "https://store77.net/apple_iphone_15_pro_max_2/",
		}
	
		var allProducts []parser.Product
	
		for site, url := range shops {
			products, err := parser.ScrapeProducts(url, site)
			if err != nil {
				log.Println("Parsing error:", site, err)
				continue
			}
			allProducts = append(allProducts, products...)
		}
	
		return c.JSON(allProducts)
	})
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
