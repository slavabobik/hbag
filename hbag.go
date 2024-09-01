package hbag

import (
	"sync"
)

// New return new instance of hbag with 0 capacity.
func New[K comparable]() *hbag[K] {
	return new[K](0)
}

// New returns new instance of hbag with specified capacity.
func NewWithCapacity[K comparable](capacity uint64) *hbag[K] {
	return new[K](capacity)
}

type hbag[K comparable] struct {
	// mu protects the state of hbag.
	mu sync.RWMutex
	// multiset stores key-count pairs, count represents the key's occurrences.
	multiset map[K]uint64
	// capacity represents the initial capacity of the multiset.
	count uint64
}

func new[K comparable](capacity uint64) *hbag[K] {
	return &hbag[K]{
		multiset: make(map[K]uint64, capacity),
		count:    0,
	}
}

// Insert adds a key to the bag and returns its previous count.
func (h *hbag[K]) Insert(key K) uint64 {
	return h.InsertMany(key, 1)
}

// InsertMany adds a key with specified count to the bag.
// If count is zero, returns count of occurencies of the key.
func (h *hbag[K]) InsertMany(key K, count uint64) uint64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	if count == 0 {
		return h.multiset[key]
	}

	last := h.multiset[key]

	h.multiset[key] += count
	h.count += count
	return last
}

// Remove decreases the count of a key in the bag by one.
// It returns the previous count of the key.
// If the key is not found, it returns 0.
func (h *hbag[K]) Remove(key K) uint64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	previous, found := h.multiset[key]
	if !found {
		return 0
	}

	switch previous {
	case 0:
	// This case should never occur as we ensure counts are always positive.
	// If it does, it indicates a bug in our implementation.
	case 1:
		delete(h.multiset, key)
	default:
		h.multiset[key]--
	}
	h.count--

	return previous
}

// Clear removes all elements from the hbag, resetting its count to zero.
// Note: This operation retains the allocated memory for the internal map.
func (h *hbag[K]) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	clear(h.multiset)
	h.count = 0
}

// Contains checks if a key exists in the bag and returns its count.
// It returns the number of occurrences of the key and a boolean indicating if the key was found.
func (h *hbag[K]) Contains(key K) (uint64, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	occurrences, found := h.multiset[key]
	return occurrences, found
}

// Get retrieves the key and its count from the bag.
// It returns the key, number of occurrences, and a boolean indicating if the key was found.
func (h *hbag[K]) Get(key K) (K, uint64, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	occurrences, found := h.multiset[key]
	return key, occurrences, found
}

// UniqLen returns the number of uniq elements in the bag.
// Duplicates are not counted.
func (h *hbag[K]) UniqLen() uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return uint64(len(h.multiset))
}

// Len returns the number of elements in the bag.
// Duplicates are counted.
func (h *hbag[K]) Len() uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.count
}

// IsUniq checks if all keys in the bag are unique.
// It returns true if each key has only one occurrence, false otherwise.
// Note: An empty bag is considered a unique bag.
func (h *hbag[K]) IsUniq() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.count == uint64(len(h.multiset))
}
