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

func ScrapeProducts(baseURL string, site string) ([]Product, error) {
	config, exists := SiteConfigs[site]
	if !exists {
		log.Println("Unknown site:", site)
		return nil, nil
	}

	c := colly.NewCollector()
	var products []Product

	// Рекурсивный вызов для сбора всех страниц
	var scrapePage func(string)
	scrapePage = func(url string) {
		c.OnHTML(config.ItemSelector, func(e *colly.HTMLElement) {
			name := e.ChildText(config.TitleSelector)
			price := CleanPrice(e.ChildText(config.PriceSelector))

			if name != "" && price != "" {
				products = append(products, Product{Name: name, Price: price, Source: config.Source})
			}
		})

		// Пагинация: ищем ссылку "Следующая страница"
		c.OnHTML(config.NextPageSelector, func(e *colly.HTMLElement) {
			nextPage := e.Attr("href")
			if nextPage != "" {
				nextURL := e.Request.AbsoluteURL(nextPage)
				log.Println("Переход на страницу:", nextURL)
				scrapePage(nextURL) // Рекурсивно парсим следующую страницу
			}
		})

		err := c.Visit(url)
		if err != nil {
			log.Println("Ошибка при посещении страницы:", err)
		}
	}

	// Запускаем парсинг с первой страницы
	scrapePage(baseURL)

	return products, nil
}
