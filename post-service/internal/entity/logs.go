package entity

import (
	"github.com/uptrace/bun"
)

type Log struct {
	bun.BaseModel `bun:"table:logs"`

	BasicEntity
	Level       *string `json:"level" bun:"level"`
	Message     *string `json:"message" bun:"message"`
	ServiceName *string `json:"service_name" bun:"service_name"`
}
