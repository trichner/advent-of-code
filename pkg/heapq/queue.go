package heapq

type node[E any] struct {
	element E
	next    *node[E]
}
type Queue[S []E, E any] struct {
	arr        S
	length     int
	comparator func(a, b E) int
}

func (q *Queue[S, E]) Push(element E) {
	p := q.length
	if len(q.arr) <= p {
		// grow
		q.grow()
	}
	q.length++

	// Insertion: Add the new element at the end of the heap, in the first available free space.
	// If this will violate the heap property, sift up the new element (swim operation)
	// until the heap property has been reestablished.
	q.arr[p] = element
	q.swim(p)
}

func (q *Queue[S, E]) swim(p int) {
	for {
		parent := (p - 1) / 2
		if q.comparator(q.arr[parent], q.arr[p]) < 0 {
			q.arr[parent], q.arr[p] = q.arr[p], q.arr[parent]
		} else {
			return
		}
		if parent == 0 {
			return
		}
		p = parent
	}
}

func (q *Queue[S, E]) grow() {
	narr := make(S, len(q.arr)*2)
	copy(narr, q.arr)
	q.arr = narr
}

func (q *Queue[S, E]) Pop() (E, bool) {
	var z E
	if q.length == 0 {
		return z, false
	}

	// Extraction: Remove the root and insert the last element of the heap in the root.
	// If this will violate the heap property, sift down the new root (sink operation) to reestablish the heap property.
	e := q.arr[0]

	q.arr[0] = q.arr[q.length-1]
	q.length--

	q.sink()

	return e, true
}

func (q *Queue[S, E]) sink() {
	p := 0
	for {
		l := 2*p + 1
		r := 2*p + 2
		if q.comparator(q.arr[l], q.arr[p]) < 0 && l < q.length {
			q.arr[l], q.arr[p] = q.arr[p], q.arr[l]
			p = l
		} else if q.comparator(q.arr[r], q.arr[p]) < 0 && r < q.length {
			q.arr[r], q.arr[p] = q.arr[p], q.arr[r]
			p = r
		} else {
			return
		}
	}
}
