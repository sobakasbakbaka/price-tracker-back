package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"price-tracker/parser"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
		DB:   0,                
	})
}

func filterAndSortProducts(products []parser.Product, c *fiber.Ctx) []parser.Product {
	searchQuery := strings.ToLower(c.Query("search"))
	shopFilter := strings.ToLower(c.Query("shop"))
	sortByPrice := c.Query("sort")

	var filtered []parser.Product
	for _, p := range products {
		nameLower := strings.ToLower(p.Name)

		if searchQuery != "" {
			words := strings.Fields(searchQuery)
			matched := true
			for _, word := range words {
				if !strings.Contains(nameLower, word) {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}
		}

		if shopFilter != "" && strings.ToLower(p.Source) != shopFilter {
			continue
		}
		filtered = append(filtered, p)
	}

	if sortByPrice == "asc" || sortByPrice == "desc" {
		sort.Slice(filtered, func(i, j int) bool {
			price1, err1 := strconv.ParseFloat(filtered[i].Price, 64)
			price2, err2 := strconv.ParseFloat(filtered[j].Price, 64)

			if err1 != nil || err2 != nil {
				return false
			}

			if sortByPrice == "asc" {
				return price1 < price2
			}
			return price1 > price2
		})
	}

	return filtered
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	initRedis()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Server is worked!"})
	})

	app.Get("/parse", func(c *fiber.Ctx) error {
		cachedData, err := redisClient.Get(context.Background(), "cached_products").Result()
		if err == nil {
			var cachedProducts []parser.Product
			if err := json.Unmarshal([]byte(cachedData), &cachedProducts); err == nil {
				return c.JSON(filterAndSortProducts(cachedProducts, c))
			}
		}
	
		shops := map[string]string{
			"indexiq": "https://indexiq.ru/catalog/iphone/",
			"biggeek": "https://biggeek.ru/catalog/apple-iphone",
			"store77": "https://store77.net/telefony_apple/",
		}
	
		var allProducts []parser.Product
	
		for site, url := range shops {
			products, err := parser.ScrapeProducts(url, site)
			if err != nil {
				log.Println("Parsing error:", site)
				continue
			}
			allProducts = append(allProducts, products...)
		}
	
		productsJSON, _ := json.Marshal(allProducts)
		redisClient.Set(context.Background(), "cached_products", productsJSON, 10*time.Minute)
	
		return c.JSON(filterAndSortProducts(allProducts, c))
	})	

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}


