package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RogerBonati/OrderApi/Model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

// generate the order key for the database

func OrderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to encode order: %w", err)
	}
	key := OrderIDKey(order.OrderId)

	txn := r.Client.TxPipeline()

	res := r.Client.SetNX(ctx, key, string(data), 0) // setNX: write order only if it does not exist; do not overwrite
	if err = res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("Failed to set the order: %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("Failed to add to orders set: %w", err)
	}
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("Failed to exec: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("Order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	key := OrderIDKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("Get order failed: %w", err)

	}
	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("Faile to decode order json: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := OrderIDKey(id)

	txn := r.Client.TxPipeline()
	err := txn.Del(ctx, key).Err()

	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("Delete order: %w", err)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		fmt.Errorf("Failed to remove from the order set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("Failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to encode order: %w", err)
	}
	key := order.OrderId(order, OrderId)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err() // SetXX do the update only if the entry exist

	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else {
		return fmt.Errorf("Set order: %w", err)
	}
	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))
	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get the order IDs: %w", err)

	}

	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get the orders: %w", err)
	}

	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order
		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("Failed to decode order json: %w", err)

		}
		orders[i] = order

	}
	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
