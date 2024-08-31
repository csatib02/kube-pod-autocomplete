package model

// Set is a generic interface for a collection of unique items
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates a new Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// Add adds an item to the set.
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// ToSlice converts the set to a slice of items.
func (s *Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.items))
	for item := range s.items {
		slice = append(slice, item)
	}
	return slice
}
