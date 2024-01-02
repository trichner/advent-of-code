package sets

import (
	"fmt"
	"strings"
)

type Set[K comparable] map[K]struct{}

func (s Set[K]) Size() int {
	return len(s)
}

func (s Set[K]) Put(e K) {
	s[e] = struct{}{}
}

func (s Set[K]) PutAll(elements []K) {
	for _, k := range elements {
		s.Put(k)
	}
}

func (s Set[K]) Has(e K) bool {
	_, ok := s[e]
	return ok
}

func (s Set[K]) Subtract(b Set[K]) {
	for k := range b {
		if s.Has(k) {
			delete(s, k)
		}
	}
}

func (s Set[K]) Keys() []K {
	var keys []K
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

func (s Set[K]) String() string {
	if len(s) == 0 {
		return "{}"
	}

	var keys []string
	for _, k := range s.Keys() {
		keys = append(keys, fmt.Sprintf("%v", k))
	}

	return "{ " + strings.Join(keys, ", ") + " }"
}

func New[K comparable]() Set[K] {
	return make(Set[K])
}

func NewFrom[K comparable](s Set[K]) Set[K] {
	n := New[K]()
	n.PutAll(s.Keys())
	return n
}

func Intersect[K comparable](a, b Set[K]) Set[K] {
	union := make(Set[K], min(len(a), len(b)))
	for k := range a {
		if b.Has(k) {
			union.Put(k)
		}
	}
	return union
}

func Union[K comparable](a, b Set[K]) Set[K] {
	union := make(Set[K], min(len(a), len(b)))
	for k := range a {
		union.Put(k)
	}
	for k := range b {
		union.Put(k)
	}
	return union
}
