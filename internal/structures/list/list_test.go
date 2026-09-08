package list

import "testing"

// checkList проверяет порядок самих узлов, границы и связи в обе стороны
func checkList[V any](t *testing.T, l List[V], want ...*Item[V]) {
	t.Helper()
	if got := l.Len(); got != len(want) {
		t.Fatalf("Len() = %d; want %d", got, len(want))
	}

	// ограничиваем обход ожидаемым числом узлов, чтобы цикл не подвесил тест
	current := l.Front()
	for index, expected := range want {
		if current != expected {
			t.Fatalf("forward node %d = %p; want %p", index, current, expected)
		}
		current = current.Next
	}
	if current != nil {
		t.Fatalf("forward traversal has extra node or cycle: %p", current)
	}

	// обратный обход проверяет Prev, а выход за голову — левую границу
	current = l.Back()
	for index := len(want) - 1; index >= 0; index-- {
		if current != want[index] {
			t.Fatalf("backward node %d = %p; want %p", index, current, want[index])
		}
		current = current.Prev
	}
	if current != nil {
		t.Fatalf("backward traversal has extra node or cycle: %p", current)
	}
}

// TestNewList проверяет размер и обе границы пустого списка
func TestNewList(t *testing.T) {
	checkList(t, NewList[string]())
}

// TestPushFront проверяет добавление в пустой и непустой список
func TestPushFront(t *testing.T) {
	l := NewList[any]()
	var want []*Item[any]
	// одинаковые значения должны создавать разные узлы, nil тоже допустим
	for _, value := range []any{"A", "A", nil, 42} {
		item := l.PushFront(value)
		if item == nil || item.Value != value {
			t.Fatalf("PushFront(%v) = %+v; want node with supplied value", value, item)
		}
		want = append([]*Item[any]{item}, want...)
		checkList(t, l, want...)
	}
}

// TestPushBack проверяет добавление в пустой и непустой список
func TestPushBack(t *testing.T) {
	l := NewList[any]()
	var want []*Item[any]
	// одинаковые значения должны создавать разные узлы, nil тоже допустим
	for _, value := range []any{"A", "A", nil, 42} {
		item := l.PushBack(value)
		if item == nil || item.Value != value {
			t.Fatalf("PushBack(%v) = %+v; want node with supplied value", value, item)
		}
		want = append(want, item)
		checkList(t, l, want...)
	}
}

// TestRemove проверяет удаление единственного узла, головы, середины и хвоста
func TestRemove(t *testing.T) {
	for _, tc := range []struct {
		name  string
		size  int
		index int // -1 означает nil
	}{
		{"nil_empty", 0, -1},
		{"nil_nonempty", 3, -1},
		{"only_node", 1, 0},
		{"head", 3, 0},
		{"middle", 3, 1},
		{"tail", 3, 2},
		{"head_of_two", 2, 0},
		{"tail_of_two", 2, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := NewList[int]()
			var want []*Item[int]
			for index := 0; index < tc.size; index++ {
				want = append(want, l.PushBack(index))
			}
			var item *Item[int]
			if tc.index >= 0 {
				// сохраняем удаляемый узел и исключаем его из ожидаемого порядка
				item = want[tc.index]
				want = append(want[:tc.index], want[tc.index+1:]...)
			}

			l.Remove(item)
			checkList(t, l, want...)
		})
	}
}

// TestMove проверяет перенос в оба конца и сохранение значений узлов
func TestMove(t *testing.T) {
	for _, operation := range []struct {
		name string
		move func(List[int], *Item[int])
	}{
		{"MoveToFront", List[int].MoveToFront},
		{"MoveToBack", List[int].MoveToBack},
	} {
		t.Run(operation.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				size  int
				index int   // -1 означает nil
				front []int // порядок исходных узлов после MoveToFront
				back  []int // порядок исходных узлов после MoveToBack
			}{
				{"nil_empty", 0, -1, nil, nil},
				{"nil_nonempty", 3, -1, []int{0, 1, 2}, []int{0, 1, 2}},
				{"only_node", 1, 0, []int{0}, []int{0}},
				{"head", 3, 0, []int{0, 1, 2}, []int{1, 2, 0}},
				{"middle", 3, 1, []int{1, 0, 2}, []int{0, 2, 1}},
				{"tail", 3, 2, []int{2, 0, 1}, []int{0, 1, 2}},
				{"head_of_two", 2, 0, []int{0, 1}, []int{1, 0}},
				{"tail_of_two", 2, 1, []int{1, 0}, []int{0, 1}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					l := NewList[int]()
					var items []*Item[int]
					for index := 0; index < tc.size; index++ {
						items = append(items, l.PushBack(index))
					}
					var item *Item[int]
					if tc.index >= 0 {
						item = items[tc.index]
					}
					order := tc.front
					if operation.name == "MoveToBack" {
						order = tc.back
					}
					var want []*Item[int]
					for _, index := range order {
						want = append(want, items[index])
					}

					operation.move(l, item)
					checkList(t, l, want...)
					// повторный перенос того же узла уже не меняет порядок
					operation.move(l, item)
					checkList(t, l, want...)
					for index, node := range items {
						if node.Value != index {
							t.Fatalf("node %d value = %v; want %d", index, node.Value, index)
						}
					}
				})
			}
		})
	}
}

// TestListSequence проверяет связи после чередования добавлений, переносов и удалений
func TestListSequence(t *testing.T) {
	l := NewList[string]()
	a := l.PushBack("A")
	checkList(t, l, a)
	b := l.PushFront("B")
	checkList(t, l, b, a)
	c := l.PushBack("C")
	checkList(t, l, b, a, c)
	l.MoveToFront(c)
	checkList(t, l, c, b, a)
	l.MoveToBack(b)
	checkList(t, l, c, a, b)
	l.Remove(a)
	checkList(t, l, c, b)
	l.Remove(c)
	checkList(t, l, b)
	l.Remove(b)
	checkList(t, l)

	// после удаления всех узлов список можно заполнить снова
	d := l.PushFront("D")
	checkList(t, l, d)
	e := l.PushBack("E")
	checkList(t, l, d, e)
}
