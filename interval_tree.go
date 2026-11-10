// Package interval_tree provides an interval tree that indexes intervals by their
// endpoints and answers stabbing queries: given a point, return every interval that
// contains it. Intervals are half-open [low, high), so an interval contains a point p
// exactly when low <= p < high. This choice makes adjacent intervals [0,10) and [10,20)
// disjoint, which is the usual convention for interval indexing and avoids double
// counting at shared endpoints.
package interval_tree

import "sort"

// Interval is a half-open interval [Low, High). Low must be less than High for the
// interval to be valid; invalid intervals are rejected by Insert.
type Interval struct {
	Low  float64
	High float64
}

// node is a single node of the balanced binary search tree. Each node stores one
// interval and the maximum High value in its subtree, which is the augmentation that
// lets stabbing queries prune whole subtrees.
type node struct {
	interval Interval
	max      float64
	left     *node
	right    *node
}

// Tree is an interval tree. The zero value is an empty tree ready for use. A Tree must
// not be copied after first use; it is not safe for concurrent use by multiple goroutines.
type Tree struct {
	root *node
	size int
}

// New returns an empty Tree.
func New() *Tree {
	return &Tree{}
}

// Len reports the number of intervals stored in the tree.
func (t *Tree) Len() int {
	return t.size
}

// Insert adds iv to the tree and returns true. It returns false and leaves the tree
// unchanged if iv is not a valid interval, i.e. if iv.Low >= iv.High. Rejecting invalid
// intervals keeps the tree's invariant (every node's High is strictly greater than its
// Low) intact and makes query results meaningful.
func (t *Tree) Insert(iv Interval) bool {
	if iv.Low >= iv.High {
		return false
	}
	t.root = insert(t.root, iv)
	t.size++
	return true
}

// insert inserts iv into the subtree rooted at n and returns the new subtree root.
// The tree is kept as an unbalanced binary search tree ordered by Low, augmented with
// the maximum High in each subtree.
func insert(n *node, iv Interval) *node {
	if n == nil {
		return &node{interval: iv, max: iv.High}
	}
	if iv.Low < n.interval.Low {
		n.left = insert(n.left, iv)
	} else {
		n.right = insert(n.right, iv)
	}
	update(n)
	return n
}

// update recomputes the max augmentation for n from its children.
func update(n *node) {
	m := n.interval.High
	if n.left != nil && n.left.max > m {
		m = n.left.max
	}
	if n.right != nil && n.right.max > m {
		m = n.right.max
	}
	n.max = m
}

// Query returns every interval in the tree that contains p. The result is sorted by
// Low ascending, then by High ascending, so the output is deterministic regardless of
// insertion order. An empty slice is returned when there are no matches; it is never nil.
func (t *Tree) Query(p float64) []Interval {
	var out []Interval
	collect(t.root, p, &out)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Low != out[j].Low {
			return out[i].Low < out[j].Low
		}
		return out[i].High < out[j].High
	})
	if out == nil {
		out = []Interval{}
	}
	return out
}

// collect walks the subtree rooted at n and appends every interval that contains p to
// out. The max augmentation lets us skip any subtree whose largest High is not greater
// than p, because no interval there can contain p.
func collect(n *node, p float64, out *[]Interval) {
	if n == nil {
		return
	}
	// A subtree can contain p only if its max High is strictly greater than p
	// (half-open intervals need High > p).
	if n.max <= p {
		return
	}
	if n.interval.Low <= p && p < n.interval.High {
		*out = append(*out, n.interval)
	}
	// Intervals in the left subtree have Low <= n.interval.Low, but their High
	// may still exceed p, so we must descend. The right subtree only needs
	// visiting if p is at least n.interval.Low.
	collect(n.left, p, out)
	if p >= n.interval.Low {
		collect(n.right, p, out)
	}
}
