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
	config, exists := SiteConfigs[site]
	if !exists {
		log.Println("Unknown site:", site)
		return nil, nil
	}

	c := colly.NewCollector()
	var products []Product

	c.OnHTML(config.ItemSelector, func(e *colly.HTMLElement) {
		name := e.ChildText(config.TitleSelector)
		price := CleanPrice(e.ChildText(config.PriceSelector))

		if name != "" && price != "" {
			products = append(products, Product{Name: name, Price: price, Source: config.Source})
		}
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return products, nil
}