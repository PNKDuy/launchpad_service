package response

type Response struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
	PriceChange string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	Volume string `json:"volume"`
}

type Currency struct {
	Symbol string `json:"symbol"`
	Price float64 `json:"price"`
	PriceChange float64 `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	Volume string `json:"volume"`
}