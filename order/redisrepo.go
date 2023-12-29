package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

var ErrNotExist = errors.New("order does not exist")

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order Order) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}
	key := orderIDKey(order.OrderID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(jsonData), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (Order, error) {
	key := orderIDKey(id)
	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return Order{}, ErrNotExist
	} else if err != nil {
		return Order{}, fmt.Errorf("get order: %w", err)
	}

	var order Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return Order{}, fmt.Errorf("failed to decode order json: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("delete order: %w", err)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order Order) error {
	jsondata, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}
	key := orderIDKey(order.OrderID)
	err = r.Client.SetXX(ctx, key, string(jsondata), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("set order: %w", err)
	}
	return nil
}

type PaginationOptions struct {
	Count  uint64
	Cursor uint64
}

type Results struct {
	Orders []Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page PaginationOptions) (Results, error) {
	keys, cursor, err := r.Client.SScan(
		ctx,
		"orders",
		page.Cursor,
		"*",
		int64(page.Count),
	).Result()

	if err != nil {
		return Results{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	if len(keys) == 0 {
		return Results{Orders: []Order{}}, nil
	}

	res, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return Results{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]Order, len(res))
	for i, v := range res {
		v := v.(string)
		var order Order
		if err := json.Unmarshal([]byte(v), &order); err != nil {
			return Results{}, fmt.Errorf("failed to decode order json: %w", err)
		}

		orders[i] = order
	}
	return Results{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
