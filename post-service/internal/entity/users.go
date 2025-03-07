package entity

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	BasicEntity
	FirstName *string `json:"first_name"     bun:"first_name"`
	LastName  *string `json:"last_name"      bun:"last_name"`
	Username  *string `json:"username"   bun:"username"`
	Email     *string `json:"email"         bun:"email"`
	Phone     *string `json:"phone"         bun:"phone"`
	Gender    *string `json:"gender"     bun:"gender"`
	Password  *string `json:"password"   bun:"password"`
	Role      *string `json:"role"       bun:"role"`
}
