package catalog

// Product represents a single product from the catalog
type Product struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
}
