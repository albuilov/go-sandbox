package lru

import (
	"fmt"
	"sync"

	"github.com/albuilov/go-sandbox/internal/structures/list"
)

// Cache — кеш с ограниченной вместимостью и вытеснением по LRU
// можно использовать из нескольких горутин; потокобезопасен
type Cache[K comparable, V any] interface {
	Len() int
	Set(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
	Clear()
}

// nodeCache хранит ключ и значение
type nodeCache[K comparable, V any] struct {
	key   K
	value V
}

type lruCache[K comparable, V any] struct {
	capacity int
	nodes    map[K]*list.Item[nodeCache[K, V]]
	queue    list.List[nodeCache[K, V]]
	mu       sync.Mutex
}

// NewCache создает пустой кеш с заданной вместимостью
// Вместимость должна быть больше нуля
func NewCache[K comparable, V any](capacity int) (Cache[K, V], error) {
	if capacity < 1 {
		return nil, fmt.Errorf("capacity must be positive, got %d", capacity)
	}

	c := lruCache[K, V]{
		capacity: capacity,
		nodes:    make(map[K]*list.Item[nodeCache[K, V]], capacity),
		queue:    list.NewList[nodeCache[K, V]](),
	}

	return &c, nil
}

// Len возвращает количество элементов в кеше
func (c *lruCache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.nodes)
}

// Set добавляет или обновляет значение и переносит элемент в начало
// При переполнении вытесняет давно использованный элемент
func (c *lruCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// ключ уже есть, обновляем значение и переносим узел в начало
	if node, exist := c.nodes[key]; exist {
		node.Value.value = value
		c.queue.MoveToFront(node)

		return
	}

	// добавляем новый узел в начало и сохраняем его в map
	newNode := c.queue.PushFront(nodeCache[K, V]{
		key:   key,
		value: value,
	})
	c.nodes[key] = newNode

	// места не хватает, удаляем хвост из списка map
	if len(c.nodes) > c.capacity {
		lastNode := c.queue.Back()

		delete(c.nodes, lastNode.Value.key)
		c.queue.Remove(lastNode)
	}
}

// Get возвращает значение и переносит найденный элемент в начало
// Если ключа нет, возвращает нулевое значение и false
func (c *lruCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exist := c.nodes[key]; exist {
		c.queue.MoveToFront(node)
		return node.Value.value, true
	}

	var zero V
	return zero, false
}

// Remove удаляет элемент по ключу
// Если ключа нет, ничего не делает
func (c *lruCache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exist := c.nodes[key]; exist {
		delete(c.nodes, key)
		c.queue.Remove(node)
	}
}

// Clear удаляет все элементы, сохраняя вместимость
func (c *lruCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nodes = make(map[K]*list.Item[nodeCache[K, V]], c.capacity)
	c.queue = list.NewList[nodeCache[K, V]]()
}
