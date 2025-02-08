package parser

type SiteConfig struct {
	ItemSelector string
	TitleSelector string
	PriceSelector string
	Source string
}

var SiteConfigs = map[string]SiteConfig{
	"indexiq": {
		ItemSelector:  ".product-item",
		TitleSelector: ".product-item__link",
		PriceSelector: ".product-item__price-visible",
		Source:        "indexiq",
	},
	"biggeek": {
		ItemSelector:  ".catalog-card",
		TitleSelector: ".catalog-card__title",
		PriceSelector: ".cart-modal-count",
		Source:        "biggeek",
	},
	"store77": {
		ItemSelector:  ".blocks_product",
		TitleSelector: ".bp_text_info",
		PriceSelector: ".bp_text_price",
		Source: 	  "store77",
	},
}
