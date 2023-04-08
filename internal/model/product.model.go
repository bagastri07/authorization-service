package model

import (
	"context"
	"fmt"
	"time"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func NewProductCacheKeyFromID(ID int) string {
	return fmt.Sprintf("product:%v", ID)
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, ID int) (*Product, error)
}

type ProductUsecase interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, ID int) (*Product, error)
}
