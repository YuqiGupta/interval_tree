package interval_tree

import (
	"reflect"
	"testing"
)

func TestEmptyTree(t *testing.T) {
	tree := New()
	if got := tree.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
	if got := tree.Query(5); len(got) != 0 {
		t.Fatalf("Query on empty tree = %v, want empty", got)
	}
}

func TestInsertRejectsInvalid(t *testing.T) {
	tree := New()
	if tree.Insert(Interval{Low: 5, High: 5}) {
		t.Fatal("Insert with Low == High should be rejected")
	}
	if tree.Insert(Interval{Low: 6, High: 5}) {
		t.Fatal("Insert with Low > High should be rejected")
	}
	if got := tree.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0 after rejected inserts", got)
	}
}

func TestInsertAndQuerySingle(t *testing.T) {
	tree := New()
	if !tree.Insert(Interval{Low: 1, High: 5}) {
		t.Fatal("Insert returned false for a valid interval")
	}
	if got := tree.Len(); got != 1 {
		t.Fatalf("Len() = %d, want 1", got)
	}
	if got := tree.Query(3); !reflect.DeepEqual(got, []Interval{{1, 5}}) {
		t.Fatalf("Query(3) = %v, want [{1 5}]", got)
	}
}

func TestHalfOpenBoundary(t *testing.T) {
	tree := New()
	tree.Insert(Interval{Low: 1, High: 5})
	// Low is inclusive.
	if got := tree.Query(1); !reflect.DeepEqual(got, []Interval{{1, 5}}) {
		t.Fatalf("Query(1) = %v, want [{1 5}]", got)
	}
	// High is exclusive.
	if got := tree.Query(5); len(got) != 0 {
		t.Fatalf("Query(5) = %v, want empty", got)
	}
}

func TestAdjacentIntervalsDoNotDoubleCount(t *testing.T) {
	tree := New()
	tree.Insert(Interval{Low: 0, High: 10})
	tree.Insert(Interval{Low: 10, High: 20})
	if got := tree.Query(10); !reflect.DeepEqual(got, []Interval{{10, 20}}) {
		t.Fatalf("Query(10) = %v, want [{10 20}]", got)
	}
}

func TestQueryReturnsAllContainingIntervals(t *testing.T) {
	tree := New()
	intervals := []Interval{
		{Low: 1, High: 5},
		{Low: 2, High: 6},
		{Low: 0, High: 10},
		{Low: 7, High: 9},
		{Low: 3, High: 4},
	}
	for _, iv := range intervals {
		tree.Insert(iv)
	}
	got := tree.Query(3)
	want := []Interval{{0, 10}, {1, 5}, {2, 6}, {3, 4}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Query(3) = %v, want %v", got, want)
	}
}

func TestQueryOutsideAll(t *testing.T) {
	tree := New()
	tree.Insert(Interval{Low: 1, High: 5})
	tree.Insert(Interval{Low: 10, High: 20})
	if got := tree.Query(7); len(got) != 0 {
		t.Fatalf("Query(7) = %v, want empty", got)
	}
}

func TestQueryDeterministicOrder(t *testing.T) {
	tree := New()
	// Insert in an order that would produce a different tree shape than sorted order.
	tree.Insert(Interval{Low: 5, High: 6})
	tree.Insert(Interval{Low: 1, High: 10})
	tree.Insert(Interval{Low: 3, High: 4})
	tree.Insert(Interval{Low: 2, High: 8})
	got := tree.Query(3)
	want := []Interval{{1, 10}, {2, 8}, {3, 4}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Query(3) = %v, want %v", got, want)
	}
}

func TestQueryReturnsNonNilEmptySlice(t *testing.T) {
	tree := New()
	got := tree.Query(42)
	if got == nil {
		t.Fatal("Query returned a nil slice; want an empty non-nil slice")
	}
	if len(got) != 0 {
		t.Fatalf("Query = %v, want empty", got)
	}
}

func TestNegativeAndFloatIntervals(t *testing.T) {
	tree := New()
	tree.Insert(Interval{Low: -2.5, High: 0})
	tree.Insert(Interval{Low: -1, High: 1.5})
	if got := tree.Query(-1.0); !reflect.DeepEqual(got, []Interval{{-2.5, 0}, {-1, 1.5}}) {
		t.Fatalf("Query(-1) = %v, want [{-2.5 0} {-1 1.5}]", got)
	}
	if got := tree.Query(0.0); !reflect.DeepEqual(got, []Interval{{-1, 1.5}}) {
		t.Fatalf("Query(0) = %v, want [{-1 1.5}]", got)
	}
}

func TestDuplicateIntervals(t *testing.T) {
	tree := New()
	tree.Insert(Interval{Low: 1, High: 5})
	tree.Insert(Interval{Low: 1, High: 5})
	got := tree.Query(3)
	want := []Interval{{1, 5}, {1, 5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Query(3) = %v, want %v", got, want)
	}
}
