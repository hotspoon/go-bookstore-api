package author

type Author struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"s_name"`
}

type Request struct {
	Name string `json:"name" binding:"required"`
}
