package model

type Post struct {
	ID            string  `json:"id"`
	AuthorID      *string `json:"authorId,omitempty"`
	Title         string  `json:"title"`
	Body          string  `json:"body"`
	AllowComments bool    `json:"allowComments"`
	CreatedAt     string  `json:"createdAt"`
	DeletedAt     *string `json:"deletedAt"`
}
