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

func TestHBag_Merge(t *testing.T) {
	tests := []struct {
		name            string
		bag1Entries     map[string]uint64
		bag2Entries     map[string]uint64
		expectedEntries map[string]uint64
		expectedLen     uint64
		expectedUniqLen uint64
	}{
		{
			name:            "Merge empty bags",
			bag1Entries:     nil,
			bag2Entries:     nil,
			expectedEntries: map[string]uint64{},
			expectedLen:     0,
			expectedUniqLen: 0,
		},
		{
			name:        "Merge empty into non-empty",
			bag1Entries: map[string]uint64{"apple": 1, "banana": 2},
			bag2Entries: nil,
			expectedEntries: map[string]uint64{
				"apple":  1,
				"banana": 2,
			},
			expectedLen:     3,
			expectedUniqLen: 2,
		},
		{
			name:        "Merge non-empty into empty",
			bag1Entries: nil,
			bag2Entries: map[string]uint64{"apple": 1, "banana": 2},
			expectedEntries: map[string]uint64{
				"apple":  1,
				"banana": 2,
			},
			expectedLen:     3,
			expectedUniqLen: 2,
		},
		{
			name: "Merge with no overlapping elements",
			bag1Entries: map[string]uint64{
				"apple":  1,
				"banana": 2,
			},
			bag2Entries: map[string]uint64{
				"cherry": 3,
				"date":   1,
			},
			expectedEntries: map[string]uint64{
				"apple":  1,
				"banana": 2,
				"cherry": 3,
				"date":   1,
			},
			expectedLen:     7,
			expectedUniqLen: 4,
		},
		{
			name: "Merge with overlapping elements",
			bag1Entries: map[string]uint64{
				"apple":  2,
				"banana": 1,
				"cherry": 3,
			},
			bag2Entries: map[string]uint64{
				"apple":  1,
				"banana": 2,
				"date":   1,
			},
			expectedEntries: map[string]uint64{
				"apple":  3,
				"banana": 3,
				"cherry": 3,
				"date":   1,
			},
			expectedLen:     10,
			expectedUniqLen: 4,
		},
		{
			name: "Merge with large counts",
			bag1Entries: map[string]uint64{
				"apple": 1000000,
			},
			bag2Entries: map[string]uint64{
				"apple": 1000000,
			},
			expectedEntries: map[string]uint64{
				"apple": 2000000,
			},
			expectedLen:     2000000,
			expectedUniqLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the first bag
			bag1 := New[string]()
			for k, v := range tt.bag1Entries {
				bag1.InsertMany(k, v)
			}

			// Initialize the second bag
			bag2 := New[string]()
			for k, v := range tt.bag2Entries {
				bag2.InsertMany(k, v)
			}

			// Perform the merge operation
			bag1.Merge(bag2)

			// Check the total length of the merged bag
			require.Equal(t, tt.expectedLen, bag1.Len(), "Total length should match expected")

			// Check the unique length of the merged bag
			require.Equal(t, tt.expectedUniqLen, bag1.UniqLen(), "Unique length should match expected")

			// Verify all expected entries are present with correct counts
			for k, expectedCount := range tt.expectedEntries {
				actualCount, found := bag1.Contains(k)
				require.True(t, found, "Expected key %s should be present", k)
				require.Equal(t, expectedCount, actualCount, "Count for key %s should match expected", k)
			}

			// Verify no unexpected entries are present
			for k, v := range bag1.multiset {
				expectedCount, exists := tt.expectedEntries[k]
				require.True(t, exists, "Unexpected key %s found in result", k)
				require.Equal(t, expectedCount, v, "Unexpected count for key %s", k)
			}

			// Verify the second bag remains unchanged
			for k, expectedCount := range tt.bag2Entries {
				actualCount, found := bag2.Contains(k)
				require.True(t, found, "Original key in second bag should still be present")
				require.Equal(t, expectedCount, actualCount, "Count in second bag should not change")
			}
		})
	}
}
