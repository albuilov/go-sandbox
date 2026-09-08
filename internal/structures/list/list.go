package list

// List — двусвязный список с добавлением и переносом узлов в оба конца
type List[V any] interface {
	Len() int
	Front() *Item[V]
	Back() *Item[V]
	PushFront(v V) *Item[V]
	PushBack(v V) *Item[V]
	Remove(i *Item[V])
	MoveToFront(i *Item[V])
	MoveToBack(i *Item[V])
}

// Item — узел двусвязного списка
// Value хранит значение того же типа, что и список
// Next указывает на следующий узел, Prev — на предыдущий
// Связи узла изменяются методами списка
type Item[V any] struct {
	Value V
	Next  *Item[V]
	Prev  *Item[V]
}

type list[V any] struct {
	head *Item[V]
	tail *Item[V]
	size int
}

// NewList создает пустой двусвязный список
func NewList[V any]() List[V] {
	return &list[V]{}
}

// Len возвращает количество узлов в списке
func (l *list[V]) Len() int {
	return l.size
}

// Front возвращает голову списка или nil, если список пуст
func (l *list[V]) Front() *Item[V] {
	return l.head
}

// Back возвращает хвост списка или nil, если список пуст
func (l *list[V]) Back() *Item[V] {
	return l.tail
}

// PushFront добавляет новое значение в начало списка и возвращает созданный узел
func (l *list[V]) PushFront(v V) *Item[V] {
	// ставим новый узел перед прежней головой
	//	nil <— head => nil <— newItem <—> head
	newItem := &Item[V]{Value: v, Next: l.head}

	if l.head != nil {
		l.head.Prev = newItem
	} else {
		// список был пуст, новый узел становится и головой, и хвостом
		l.tail = newItem
	}

	l.head = newItem
	l.size++

	return newItem
}

// PushBack добавляет новое значение в конец списка и возвращает созданный узел
func (l *list[V]) PushBack(v V) *Item[V] {
	// ставим новый узел после прежнего хвоста
	//	tail —> nil => tail <—> newItem —> nil
	newItem := &Item[V]{Value: v, Prev: l.tail}

	if l.tail != nil {
		l.tail.Next = newItem
	} else {
		// список был пуст, новый узел становится и головой, и хвостом
		l.head = newItem
	}

	l.tail = newItem
	l.size++

	return newItem
}

// Remove удаляет узел из списка
// Ненулевой i должен принадлежать этому списку
func (l *list[V]) Remove(i *Item[V]) {
	// передан nil или список пуст
	if i == nil || l.size == 0 {
		return
	}

	// соединяем соседей в обход i
	//	prev <—> i <—> next => prev <—> next

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		// i был головой, теперь головой становится следующий узел
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		// i был хвостом, теперь хвостом становится предыдущий узел
		l.tail = i.Prev
	}

	l.size--
}

// MoveToFront переносит существующий узел в начало списка
// Ненулевой i должен принадлежать этому списку
func (l *list[V]) MoveToFront(i *Item[V]) {
	// передан nil или узел уже является головой
	if i == nil || l.head == i {
		return
	}

	prev, next := i.Prev, i.Next

	// отсоединяем i
	//	prev <—> i <—> next => prev <—> next
	// prev существует, так как i не является головой
	prev.Next = next
	if next != nil {
		next.Prev = prev
	} else {
		// i был хвостом, теперь новым хвостом становится prev
		l.tail = prev
	}

	// ставим i перед прежней головой
	// 	nil <— head => nil <— i <—> head
	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i

	// размер списка не меняется
}

// MoveToBack переносит существующий узел в конец списка
// Ненулевой i должен принадлежать этому списку
func (l *list[V]) MoveToBack(i *Item[V]) {
	// передан nil или узел уже является хвостом
	if i == nil || l.tail == i {
		return
	}

	prev, next := i.Prev, i.Next

	// отсоединяем i
	//	prev <—> i <—> next => prev <—> next
	// next существует, так как i не является хвостом
	next.Prev = prev
	if prev != nil {
		prev.Next = next
	} else {
		// i был головой, теперь новой головой становится next
		l.head = next
	}

	// ставим i после прежнего хвоста
	//	tail —> nil => tail <—> i —> nil
	i.Next = nil
	i.Prev = l.tail
	l.tail.Next = i
	l.tail = i

	// размер списка не меняется
}
