package ttl

import (
	"fmt"
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Len() int
	Set(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
	Clear()
	DeleteExpired()
}

type nodeCache[V any] struct {
	value     V
	expiresAt time.Time
}

type ttlCache[K comparable, V any] struct {
	ttl   time.Duration
	nodes map[K]nodeCache[V]
	now   func() time.Time
	mu    sync.Mutex
}

func NewCache[K comparable, V any](ttl time.Duration) (Cache[K, V], error) {
	if ttl <= 0 {
		return nil, fmt.Errorf("ttl must be positive, got %s", ttl)
	}

	c := ttlCache[K, V]{
		ttl:   ttl,
		nodes: make(map[K]nodeCache[V]),
		now:   time.Now,
	}

	return &c, nil
}

// Len удаляет просроченные записи и возвращает количество активных
func (c *ttlCache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.deleteExpired()

	return len(c.nodes)
}

// Set сохраняет значение и начинает срок жизни заново
func (c *ttlCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nodes[key] = nodeCache[V]{
		value:     value,
		expiresAt: c.now().Add(c.ttl),
	}
}

// Get возвращает значение, если срок жизни еще не истек
// Чтение не продлевает TTL
func (c *ttlCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exist := c.nodes[key]; exist {
		if c.now().Before(node.expiresAt) {
			return node.value, true
		}

		// срок истек, удаляем запись из map
		delete(c.nodes, key)
	}

	var zero V
	return zero, false
}

// Remove удаляет элемент по ключу
// Если ключа нет, ничего не делает
func (c *ttlCache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exist := c.nodes[key]; exist {
		delete(c.nodes, key)
	}
}

// Clear удаляет все записи, сохраняя настройку TTL
func (c *ttlCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nodes = make(map[K]nodeCache[V])
}

// DeleteExpired удаляет записи, срок которых истек к моменту now
func (c *ttlCache[K, V]) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.deleteExpired()
}

// deleteExpired удаляет записи, срок которых истек к моменту now
// вызываем только под блокировкой mu
func (c *ttlCache[K, V]) deleteExpired() {
	now := c.now()

	for key, node := range c.nodes {
		if !now.Before(node.expiresAt) {
			delete(c.nodes, key)
		}
	}
}
