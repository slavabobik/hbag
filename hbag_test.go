package hbag

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHBag_Insert(t *testing.T) {
	tests := []struct {
		name            string
		initialEntries  map[string]uint64
		keyToInsert     string
		expectedCount   uint64
		expectedLen     uint64
		expectedUniqLen uint64
	}{
		{
			name:            "Insert into empty bag",
			initialEntries:  nil,
			keyToInsert:     "apple",
			expectedCount:   0,
			expectedLen:     1,
			expectedUniqLen: 1,
		},
		{
			name: "Insert new key",
			initialEntries: map[string]uint64{
				"banana": 2,
				"cherry": 1,
			},
			keyToInsert:     "apple",
			expectedCount:   0,
			expectedLen:     4,
			expectedUniqLen: 3,
		},
		{
			name: "Insert existing key",
			initialEntries: map[string]uint64{
				"apple":  1,
				"banana": 2,
			},
			keyToInsert:     "apple",
			expectedCount:   1,
			expectedLen:     4,
			expectedUniqLen: 2,
		},
		{
			name: "Insert key with high count",
			initialEntries: map[string]uint64{
				"apple": 1000000,
			},
			keyToInsert:     "apple",
			expectedCount:   1000000,
			expectedLen:     1000001,
			expectedUniqLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the hbag with the initial entries
			h := New[string]()
			for k, v := range tt.initialEntries {
				h.InsertMany(k, v)
			}

			// Perform the insert operation
			prevCount := h.Insert(tt.keyToInsert)

			// Check the returned previous count
			require.Equal(t, tt.expectedCount, prevCount, "Previous count should match")

			// Check the new count for the inserted key
			newCount, found := h.Contains(tt.keyToInsert)
			require.True(t, found, "Key should be found after insertion")
			require.Equal(t, tt.expectedCount+1, newCount, "New count should be previous count + 1")

			// Check the total length of the bag
			require.Equal(t, tt.expectedLen, h.Len(), "Total length should match expected")

			// Check the unique length of the bag
			require.Equal(t, tt.expectedUniqLen, h.UniqLen(), "Unique length should match expected")

			// Verify that other entries were not affected
			for k, v := range tt.initialEntries {
				if k != tt.keyToInsert {
					count, found := h.Contains(k)
					require.True(t, found, "Original key should still be present")
					require.Equal(t, v, count, "Count for original key should not change")
				}
			}
		})
	}
}

func TestHBag_InsertMany(t *testing.T) {
	tests := []struct {
		name            string
		initialEntries  map[string]uint64
		keyToInsert     string
		countToInsert   uint64
		expectedCount   uint64
		expectedLen     uint64
		expectedUniqLen uint64
	}{
		{
			name:            "Insert multiple into empty bag",
			initialEntries:  nil,
			keyToInsert:     "apple",
			countToInsert:   5,
			expectedCount:   0,
			expectedLen:     5,
			expectedUniqLen: 1,
		},
		{
			name: "Insert multiple of new key",
			initialEntries: map[string]uint64{
				"banana": 2,
				"cherry": 1,
			},
			keyToInsert:     "apple",
			countToInsert:   3,
			expectedCount:   0,
			expectedLen:     6,
			expectedUniqLen: 3,
		},
		{
			name: "Insert multiple of existing key",
			initialEntries: map[string]uint64{
				"apple":  2,
				"banana": 1,
			},
			keyToInsert:     "apple",
			countToInsert:   3,
			expectedCount:   2,
			expectedLen:     6,
			expectedUniqLen: 2,
		},
		{
			name: "Insert zero count",
			initialEntries: map[string]uint64{
				"apple": 5,
			},
			keyToInsert:     "apple",
			countToInsert:   0,
			expectedCount:   5,
			expectedLen:     5,
			expectedUniqLen: 1,
		},
		{
			name: "Insert large count",
			initialEntries: map[string]uint64{
				"apple": 1000000,
			},
			keyToInsert:     "apple",
			countToInsert:   1000000,
			expectedCount:   1000000,
			expectedLen:     2000000,
			expectedUniqLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the hbag with the initial entries
			h := New[string]()
			for k, v := range tt.initialEntries {
				h.InsertMany(k, v)
			}

			// Perform the InsertMany operation
			prevCount := h.InsertMany(tt.keyToInsert, tt.countToInsert)

			// Check the returned previous count
			require.Equal(t, tt.expectedCount, prevCount, "Previous count should match")

			// Check the new count for the inserted key
			newCount, found := h.Contains(tt.keyToInsert)
			require.True(t, found, "Key should be found after insertion")
			if tt.countToInsert > 0 {
				require.Equal(t, tt.expectedCount+tt.countToInsert, newCount, "New count should be previous count + inserted count")
			} else {
				require.Equal(t, tt.expectedCount, newCount, "Count should not change for zero insertion")
			}

			// Check the total length of the bag
			require.Equal(t, tt.expectedLen, h.Len(), "Total length should match expected")

			// Check the unique length of the bag
			require.Equal(t, tt.expectedUniqLen, h.UniqLen(), "Unique length should match expected")

			// Verify that other entries were not affected
			for k, v := range tt.initialEntries {
				if k != tt.keyToInsert {
					count, found := h.Contains(k)
					require.True(t, found, "Original key should still be present")
					require.Equal(t, v, count, "Count for original key should not change")
				}
			}
		})
	}
}

func TestHBag_IsUniq(t *testing.T) {
	tests := []struct {
		name           string
		entries        map[string]uint64
		expectedIsUniq bool
	}{
		{
			name:           "Empty bag",
			entries:        nil,
			expectedIsUniq: true,
		},
		{
			name: "Single entry",
			entries: map[string]uint64{
				"apple": 1,
			},
			expectedIsUniq: true,
		},
		{
			name: "Multiple unique entries",
			entries: map[string]uint64{
				"apple":  1,
				"banana": 1,
				"cherry": 1,
			},
			expectedIsUniq: true,
		},
		{
			name: "One non-unique entry",
			entries: map[string]uint64{
				"apple":  2,
				"banana": 1,
				"cherry": 1,
			},
			expectedIsUniq: false,
		},
		{
			name: "All non-unique entries",
			entries: map[string]uint64{
				"apple":  2,
				"banana": 3,
				"cherry": 4,
			},
			expectedIsUniq: false,
		},
		{
			name: "Large count for one entry",
			entries: map[string]uint64{
				"apple": 1000000,
			},
			expectedIsUniq: false,
		},
		{
			name: "Mix of unique and non-unique",
			entries: map[string]uint64{
				"apple":  1,
				"banana": 2,
				"cherry": 1,
				"date":   3,
			},
			expectedIsUniq: false,
		},
		{
			name: "Edge case: overflow attempt",
			entries: map[string]uint64{
				"apple": ^uint64(0), // Max uint64 value
			},
			expectedIsUniq: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New[string]()
			for k, v := range tt.entries {
				h.InsertMany(k, v)
			}

			isUniq := h.IsUniq()
			require.Equal(t, tt.expectedIsUniq, isUniq, "IsUniq result should match expected")

			// Additional checks
			if tt.expectedIsUniq {
				require.Equal(t, h.Len(), h.UniqLen(), "For unique bags, Len and UniqLen should be equal")
			} else {
				require.Greater(t, h.Len(), h.UniqLen(), "For non-unique bags, Len should be greater than UniqLen")
			}

			// Check consistency with individual counts
			for k, v := range tt.entries {
				count, found := h.Contains(k)
				if v > 0 {
					require.True(t, found, "Key should be found")
					require.Equal(t, v, count, "Count should match")
					if v > 1 {
						require.False(t, isUniq, "Bag should not be unique if any count is greater than 1")
					}
				} else {
					require.False(t, found, "Key with zero count should not be found")
				}
			}
		})
	}
}

// import (
// 	"sync"
// 	"testing"

// "github.com/stretchr/testify/require"
// )

// func TestNew(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[int]()

// 	require.Zero(hb.Len())
// 	require.Zero(hb.UniqLen())
// }

// func TestNewWithCapacity(t *testing.T) {
// 	require := require.New(t)
// 	hb := NewWithCapacity[int](10)

// 	require.Zero(hb.Len())
// 	require.Zero(hb.UniqLen())
// }

// func TestInsert(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[string]()

// 	last := hb.Insert("navigation")
// 	require.Equal(uint64(0), last)

// 	last = hb.Insert("navigation")
// 	require.Equal(uint64(1), last)

// 	require.Equal(uint64(2), hb.Len())
// 	require.Equal(uint64(1), hb.UniqLen())

// 	hb.Insert("meteorology")
// 	hb.Insert("meteorology")

// 	require.Equal(uint64(4), hb.Len())
// 	require.Equal(uint64(2), hb.UniqLen())
// }

// func TestInsertConcurrent(t *testing.T) {

// 	// TODO t.Parallel()? in all tests?
// 	require := require.New(t)
// 	hb := New[string]()

// 	goroutinesCount := 10_000
// 	var wg sync.WaitGroup
// 	for i := 0; i < goroutinesCount; i++ {
// 		wg.Add(1)
// 		go func(w *sync.WaitGroup) {
// 			defer wg.Done()
// 			hb.Insert("test")
// 		}(&wg)
// 	}
// 	wg.Wait()

// 	require.Equal(uint64(goroutinesCount), hb.Len())
// 	require.Equal(uint64(1), hb.UniqLen())
// }

// func TestInsertMany(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[string]()

// 	previousCount := hb.InsertMany("aerodynamics", 5)
// 	require.Equal(uint64(0), previousCount)

// 	previousCount = hb.InsertMany("aerodynamics", 4)
// 	require.Equal(uint64(5), previousCount)
// }

// func TestInsertManyZeroCount(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[string]()

// 	previousCount := hb.InsertMany("subzero", 0)
// 	require.Equal(uint64(0), previousCount)

// 	hb.InsertMany("subzero", 5)

// 	previousCount = hb.InsertMany("subzero", 0)
// 	require.Equal(uint64(5), previousCount)
// 	require.EqualValues(5, previousCount)
// }

// func TestRemove(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[string]()

// 	occ := hb.Remove("not_exists")
// 	require.EqualValues(0, occ)

// 	occ = hb.Insert("remove_me")
// 	require.EqualValues(0, occ)

// 	occ = hb.Remove("remove_me")
// 	require.EqualValues(1, occ)
// 	require.EqualValues(0, hb.Len())
// 	require.EqualValues(0, hb.UniqLen())

// 	require.EqualValues(0, hb.count)
// 	require.NotContains(hb.multiset, "remove_me")

// 	occ = hb.Insert("foo")
// 	occ = hb.Insert("foo")
// 	require.EqualValues(1, occ)

// 	occ = hb.Remove("foo")
// 	require.EqualValues(2, occ) // TODO: is it possible to extact this 4 checks into separate func?
// 	require.EqualValues(1, hb.count)
// 	require.EqualValues(1, hb.Len())
// 	require.EqualValues(1, hb.UniqLen())

// }

// func TestClear(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[int]()

// 	hb.Insert(1)
// 	hb.Insert(2)

// 	hb.Clear()

// 	require.Len(hb.multiset, 0)
// 	require.Zero(hb.count)
// }

// func TestIsUniq(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[int]()

// 	hb.Insert(1)
// 	hb.Insert(2)
// 	hb.Insert(3)

// 	require.True(hb.IsUniq())

// 	hb.Insert(1)
// 	require.False(hb.IsUniq())
// }

// func TestContains(t *testing.T) {
// 	require := require.New(t)
// 	hb := New[string]()

// 	occ, found := hb.Contains("not_exists")
// 	require.EqualValues(0, occ)
// 	require.False(found)

// 	hb.Insert("key")

// 	occ, found = hb.Contains("key")
// 	require.EqualValues(1, occ)
// 	require.True(found)
// }

// // func TestInsertManyOverflow(t *testing.T) {
// // 	require := require.New(t)

// // 	hb := New[string]()
// // 	hb.InsertMany("human_performance", math.MaxUint64)
// // 	hb.InsertMany("human_performance", 1)

// // }
