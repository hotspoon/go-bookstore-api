package book

// Book represents one row in the Book table.
type Book struct {
	ID       string  `json:"id" db:"id"`
	PubID    int     `json:"pub_id" db:"pub_id"`
	Title    string  `json:"title" db:"title"`
	Price    float64 `json:"price" db:"price"`
	Category *string `json:"category,omitempty" db:"category"`
	Quantity int     `json:"quantity" db:"quantity"`
	Format   *string `json:"b_format,omitempty" db:"b_format"`
	ProdYear int     `json:"prod_year" db:"prod_year"`
	FileSize *int    `json:"filesize,omitempty" db:"filesize"`
}
