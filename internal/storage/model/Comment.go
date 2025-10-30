package model

type Comment struct {
	ID              string  `json:"id"`
	AuthorID        *string `json:"authorId,omitempty"`
	Body            string  `json:"body"`
	ParentCommentID *string `json:"parentCommentId,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	DeletedAt       string  `json:"deletedAt"`
}
