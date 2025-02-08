package parser

import (
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	Name  string
	Price string
}

func CleanPrice(price string) string {
	price = strings.TrimSpace(price)

	re := regexp.MustCompile(`[^\d.]`)
	price = re.ReplaceAllString(price, "")

	return price
}

func ScrapeProducts(url string) ([]Product, error) {
	c := colly.NewCollector()

	var products []Product

	c.OnHTML(".product-item", func(e *colly.HTMLElement) {
		name := e.ChildText(".product-item__link")
		price := e.ChildText(".product-item__price-visible")
		price = CleanPrice(price)

		if name != "" && price != "" {
			products = append(products, Product{Name: name, Price: price})
		}
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Ошибка парсинга:", err)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return products, nil
}