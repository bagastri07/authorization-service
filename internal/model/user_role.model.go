package model

import (
	"context"

	"github.com/google/uuid"
)

type UserRole struct {
	ID     uuid.UUID `gorm:"default:uuid_generate_v4()"`
	RoleID uuid.UUID
	UserID uuid.UUID
}

type UserRoleRepository interface {
	Create(ctx context.Context, userRole *UserRole) error
}
