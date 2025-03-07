package entity

import (
	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:posts"`

	BasicEntity
	UserID  *int64  `json:"user_id"     bun:"user_id"`
	Title   *string `json:"title"      bun:"title"`
	Content *string `json:"content"   bun:"content"`
}
