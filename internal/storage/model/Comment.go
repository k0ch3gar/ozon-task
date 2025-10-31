package model

type Comment struct {
	ID              string  `json:"id"`
	AuthorID        *string `json:"authorId,omitempty"`
	Body            string  `json:"body"`
	ParentPostID    string  `json:"parentPostId"`
	ParentCommentID *string `json:"parentCommentId,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	DeletedAt       *string `json:"deletedAt"`
}
