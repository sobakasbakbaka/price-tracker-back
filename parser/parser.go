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
	Source string
}

func CleanPrice(price string) string {
	price = strings.TrimSpace(price)

	re := regexp.MustCompile(`[^\d.]`)
	price = re.ReplaceAllString(price, "")

	return price
}

func ScrapeProducts(url string, site string) ([]Product, error) {
	c := colly.NewCollector()
	var products []Product

	switch site {
		case "indexiq":
			c.OnHTML(".product-item", func(e *colly.HTMLElement) {
				name := e.ChildText(".product-item__link")
				price := e.ChildText(".product-item__price-visible")
				price = CleanPrice(price)
		
				if name != "" && price != "" {
					products = append(products, Product{Name: name, Price: price, Source: "indexiq"})
				}
			})
		case "biggeek":
			c.OnHTML(".catalog-card", func(e *colly.HTMLElement) {
				name := e.ChildText(".catalog-card__title")
				price := e.ChildText(".cart-modal-count")
				price = CleanPrice(price)
		
				if name != "" && price != "" {
					products = append(products, Product{Name: name, Price: price, Source: "biggeek"})
				}
			})
		default:
			log.Println("Unknown site")
			return nil, nil
	}

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Ошибка парсинга:", err)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return products, nil
}