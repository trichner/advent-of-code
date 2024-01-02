package queue

type node[E any] struct {
	element E
	next    *node[E]
}
type Queue[E any] struct {
	head, tail *node[E]
}

func (q *Queue[E]) Push(element E) {
	n := &node[E]{element: element}
	if q.head == nil {
		q.head = n
		q.tail = q.head
		return
	}

	q.tail.next = n
	q.tail = n
}

func (q *Queue[E]) Pop() (E, bool) {
	var z E
	if q.head == nil {
		return z, false
	}
	e := q.head
	q.head = q.head.next
	return e.element, true
}

func (q *Queue[E]) ForEach(consumer func(e E)) {
	current := q.head
	for {
		if current == nil {
			return
		}
		consumer(current.element)
		current = current.next
	}
}
