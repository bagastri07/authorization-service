package model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"default:uuid_generate_v4()"`
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserCacheKeyFromID(ID string) string {
	return fmt.Sprintf("user:%s", ID)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*uuid.UUID, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, user *User) (*TokenResp, error)

	InjectUserRepository(userRepo UserRepository)
	InjectUserRoleRepository(userRoleRepo UserRoleRepository)
}
