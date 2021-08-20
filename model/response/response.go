package response

type Response struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
	PriceChange string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	Volume string `json:"volume"`
}