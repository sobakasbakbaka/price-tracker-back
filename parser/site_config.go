package parser

type SiteConfig struct {
	ItemSelector string
	TitleSelector string
	PriceSelector string
	NextPageSelector string
	Source string
}

var SiteConfigs = map[string]SiteConfig{
	"indexiq": {
		ItemSelector:  ".product-item",
		TitleSelector: ".product-item__link",
		PriceSelector: ".product-item__price-visible",
		NextPageSelector: ".rs-pagination-more",
		Source:        "indexiq",
	},
	"biggeek": {
		ItemSelector:  ".catalog-card",
		TitleSelector: ".catalog-card__title",
		PriceSelector: ".cart-modal-count",
		NextPageSelector: ".prod-pagination__item-next",
		Source:        "biggeek",
	},
	"store77": {
		ItemSelector:  ".blocks_product",
		TitleSelector: ".bp_text_info",
		PriceSelector: ".bp_text_price",
		NextPageSelector: ".pagin_arrow .pag_right",
		Source: 	  "store77",
	},
}
