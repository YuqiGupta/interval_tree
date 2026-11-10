# interval_tree

Stores half-open intervals `[Low, High)` and answers stabbing queries: given a point `p`, return every stored interval that contains it.

## Usage

```go
package main

import (
	"fmt"

	"interval_tree"
)

func main() {
	tree := interval_tree.New()
	tree.Insert(interval_tree.Interval{Low: 1, High: 5})
	tree.Insert(interval_tree.Interval{Low: 2, High: 6})
	tree.Insert(interval_tree.Interval{Low: 0, High: 10})

	fmt.Println(tree.Len())        // 3
	fmt.Println(tree.Query(3))     // [{0 10} {1 5} {2 6}]
	fmt.Println(tree.Query(7))     // []
}
```

`Insert` returns `false` and does nothing when `Low >= High`. `Query` always returns a non-nil slice sorted by `Low` ascending, then `High` ascending.

## Why this exists

A linear scan over all intervals is the obvious way to answer a stabbing query, but it costs `O(n)` per query. This tree augments a binary search tree ordered by `Low` with the maximum `High` in each subtree. A subtree whose maximum `High` is not greater than `p` cannot contain `p`, so the query prunes it. The trade-off is that the tree is a plain unbalanced BST: insertion order can degrade it to `O(n)` depth. Keeping the code short and dependency-free was chosen over the extra complexity of balancing, since the intended workload is build-once, query-many with modest sizes.

## The awkward edge

Intervals are half-open. An interval `[0, 10)` contains `0` but not `10`, so `[0, 10)` and `[10, 20)` never both match the same point. If you want closed intervals, shift the upper endpoint outward before inserting, and remember that a point exactly on a shared boundary will then match both intervals.
