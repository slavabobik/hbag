# HBag: An implementation of a [bag](https://en.wikipedia.org/wiki/Set_(abstract_data_type)#Multiset), allowing duplicates and tracking their counts.
Inspired by [Jon Gjengset's](https://github.com/jonhoo) [hashbag](https://github.com/jonhoo/hashbag) implementation in Rust.

![build workflow](https://github.com/slavabobik/hbag/actions/workflows/go.yml/badge.svg)

## Installation

```bash
 go get github.com/slavabobik/hbag
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/slavabobik/hbag"
)

func main() {
    // Create a new bag
    bag := hbag.New[string]()

    // Insert single items
    bag.Insert("apple")              // adds one apple
    bag.InsertMany("banana", 3)      // adds three bananas
    
    // Check if items exist and get their counts
    count, exists := bag.Contains("apple")
    fmt.Printf("Apple count: %d, exists: %v\n", count, exists)
    // Output: Apple count: 1, exists: true
    
    count, exists = bag.Contains("banana")
    fmt.Printf("Banana count: %d, exists: %v\n", count, exists)
    // Output: Banana count: 3, exists: true
    
    // Get total number of items (including duplicates)
    fmt.Printf("Total items: %d\n", bag.Len())
    // Output: Total items: 4
    
    // Get number of unique items
    fmt.Printf("Unique items: %d\n", bag.UniqLen())
    // Output: Unique items: 2
    
    // Check if all items are unique
    fmt.Printf("Is unique: %v\n", bag.IsUniq())
    // Output: Is unique: false
    
    // Remove an item
    prevCount := bag.Remove("banana")
    fmt.Printf("Previous banana count: %d\n", prevCount)
    // Output: Previous banana count: 3
    
    // Clear all items
    bag.Clear()
    fmt.Printf("Items after clear: %d\n", bag.Len())
    // Output: Items after clear: 0
}
```

The HBag is thread-safe and can be safely used from multiple goroutines. You can also create a bag with initial capacity:

```go
// Create a bag with initial capacity of 100 items
bag := hbag.NewWithCapacity[int](100)
```
