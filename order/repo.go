package order

import "context"

type Repo interface {
	Insert(ctx context.Context, order Order) error
	FindByID(ctx context.Context, id uint64) (Order, error)
	DeleteByID(ctx context.Context, id uint64) error
	Update(ctx context.Context, order Order) error
	FindAll(ctx context.Context, page PaginationOptions) (Results, error)
}

type PaginationOptions struct {
	Count  uint64
	Cursor uint64
}

type Results struct {
	Orders []Order
	Cursor uint64
}
