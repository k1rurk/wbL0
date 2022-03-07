package cache

import (
	"errors"
	"sync"
	"wb_l0/database"
)

type Cache struct {
	mu     sync.Mutex
	orders map[string]database.Order
}

func NewCache(orders []database.Order) *Cache {
	o := make(map[string]database.Order)

	for i := range orders {
		o[orders[i].OrderUid] = orders[i]
	}
	return &Cache{
		orders: o,
	}
}

func (cache *Cache) Get(key string) (database.Order, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	order, ok := cache.orders[key]
	if !ok {
		err := errors.New("Key is not found!")
		return database.Order{}, err
	}

	return order, nil
}

func (cache *Cache) Set(key string, order *database.Order) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.orders[key] = *order
}

func (cache *Cache) Check(key string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if _, ok := cache.orders[key]; ok {
		return errors.New("Order already exists!")
	}
	return nil
}
