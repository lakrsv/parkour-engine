package engine

type SparseSetEntry[T comparable] struct {
	id    uint32
	value T
}
type SparseSet[T comparable] struct {
	dense    []SparseSetEntry[T]
	sparse   []uint32
	capacity uint32
	n        uint32
	null     uint32
}

func NewSparseSet[T comparable](capacity uint32) *SparseSet[T] {
	dense := make([]SparseSetEntry[T], capacity+1)
	sparse := make([]uint32, capacity+1)
	null := capacity
	for idx := range sparse {
		sparse[idx] = null
	}
	dense[null] = SparseSetEntry[T]{id: null}
	return &SparseSet[T]{dense: dense, sparse: sparse, capacity: capacity, n: 0, null: null}
}

func (set *SparseSet[T]) Insert(id uint32, value T) bool {
	if set.n >= set.capacity {
		return false
	}
	if id >= set.null {
		return false
	}
	if set.Contains(id) {
		return false
	}
	set.dense[set.n] = SparseSetEntry[T]{id: id, value: value}
	set.sparse[id] = set.n
	set.n++
	return true
}

func (set *SparseSet[T]) Remove(id uint32) bool {
	if set.n == 0 {
		return false
	}
	if set.Contains(id) {
		tmp := set.dense[set.n-1]
		set.dense[set.sparse[id]] = tmp
		set.sparse[tmp.id] = set.sparse[id]
		set.sparse[id] = set.null
		set.n--
		return true
	}
	return false
}

func (set *SparseSet[T]) Contains(id uint32) bool {
	if id >= set.null {
		return false
	}
	return set.sparse[id] != set.null
}

func (set *SparseSet[T]) Clear() {
	set.n = 0
	for idx := range set.sparse {
		set.sparse[idx] = set.null
	}
}

func (set *SparseSet[T]) Get(id uint32) (T, bool) {
	if !set.Contains(id) {
		return set.dense[set.null].value, false
	}
	return set.dense[set.sparse[id]].value, true
}

func (set *SparseSet[T]) IntersectId(other *SparseSet[T]) *SparseSet[T] {
	capacity := min(set.capacity, other.capacity)
	result := NewSparseSet[T](capacity)

	if set.n > other.n {
		smallIter := other.Iterator()
		for {
			id, item, ok := smallIter.Next()
			if !ok {
				break
			}
			if set.Contains(id) {
				result.Insert(id, item)
			}
		}
	} else {
		smallIter := set.Iterator()
		for {
			id, item, ok := smallIter.Next()
			if !ok {
				break
			}
			if other.Contains(id) {
				result.Insert(id, item)
			}
		}
	}
	return result
}

func (set *SparseSet[T]) UnionId(other *SparseSet[T]) *SparseSet[T] {
	// TODO: How to be more clever here?
	//capacity := set.n + other.n
	capacity := max(set.capacity, other.capacity)
	result := NewSparseSet[T](capacity)
	setIter := set.Iterator()
	for {
		id, item, ok := setIter.Next()
		if !ok {
			break
		}
		result.Insert(id, item)
	}
	otherIter := other.Iterator()
	for {
		id, item, ok := otherIter.Next()
		if !ok {
			break
		}
		result.Insert(id, item)
	}
	return result
}

func (set *SparseSet[T]) DifferenceId(other *SparseSet[T]) *SparseSet[T] {
	capacity := set.capacity
	result := NewSparseSet[T](capacity)

	iterator := set.Iterator()
	for {
		id, item, ok := iterator.Next()
		if !ok {
			break
		}
		if other.Contains(id) {
			continue
		}
		result.Insert(id, item)
	}
	return result
}

func (set *SparseSet[T]) Iterator() *SparseSetIterator[T] {
	return &SparseSetIterator[T]{set: set}
}

func (set *SparseSet[T]) IsEmpty() bool {
	return set.n == 0
}

func (set *SparseSet[T]) Len() uint32 {
	return set.n
}

func (set *SparseSet[T]) CopyId() *SparseSet[uint32] {
	result := NewSparseSet[uint32](set.capacity)
	iterator := set.Iterator()
	for {
		id, _, ok := iterator.Next()
		if !ok {
			break
		}
		result.Insert(id, id)
	}
	return result
}

type SparseSetIterator[T comparable] struct {
	set *SparseSet[T]
	idx uint32
}

func (iterator *SparseSetIterator[T]) Next() (id uint32, value T, ok bool) {
	if iterator.idx >= iterator.set.n {
		null := iterator.set.dense[iterator.set.null]
		return null.id, null.value, false
	}
	next := iterator.set.dense[iterator.idx]
	iterator.idx++
	//slog.Info(string(iterator.idx), "stack", getStack())
	return next.id, next.value, true
}
