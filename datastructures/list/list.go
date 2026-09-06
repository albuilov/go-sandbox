package list

type List interface {
	Len() int
	Front() *Item
	Back() *Item
	PushFront(v any) *Item
	PushBack(v any) *Item
	Remove(i *Item)
	MoveToFront(i *Item)
	MoveToBack(i *Item)
}

// Item — узел двусвязного списка
// Next указывает на следующий узел, Prev — на предыдущий
// связи узла изменяются методами списка
type Item struct {
	Value any
	Next  *Item
	Prev  *Item
}

type list struct {
	head *Item
	tail *Item
	size int
}

// NewList создаёт пустой двусвязный список
func NewList() List {
	return &list{}
}

// Len возвращает количество узлов в списке
func (l *list) Len() int {
	return l.size
}

// Front возвращает голову списка или nil, если список пуст
func (l *list) Front() *Item {
	return l.head
}

// Back возвращает хвост списка или nil, если список пуст
func (l *list) Back() *Item {
	return l.tail
}

// PushFront добавляет новое значение в начало списка и возвращает созданный узел
func (l *list) PushFront(v any) *Item {
	// ставим новый узел перед прежней головой
	//	nil <— head => nil <— newItem <—> head
	newItem := &Item{Value: v, Next: l.head}

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
func (l *list) PushBack(v any) *Item {
	// ставим новый узел после прежнего хвоста
	//	tail —> nil => tail <—> newItem —> nil
	newItem := &Item{Value: v, Prev: l.tail}

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
// ненулевой i должен принадлежать этому списку
func (l *list) Remove(i *Item) {
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
// ненулевой i должен принадлежать этому списку
func (l *list) MoveToFront(i *Item) {
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
// ненулевой i должен принадлежать этому списку
func (l *list) MoveToBack(i *Item) {
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
