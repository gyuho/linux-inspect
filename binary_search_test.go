package psn

import (
	"reflect"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	nums := []int64{1, 2, 3, 4, 5}
	if idx := BinarySearchInt64(nums, 3); idx != 2 {
		t.Fatalf("expected 2, got %d", idx)
	}
	if idx := BinarySearchInt64(nums, 6); idx != -1 {
		t.Fatalf("expected -1, got %d", idx)
	}
}

func TestBinarySearchFloat64Node(t *testing.T) {
	root := newFloat64Node(0, 5.22)
	insert(root, 1, 4.85)
	insert(root, 2, 5.77)
	insert(root, 3, 7.999)

	if root.Left.Value != 4.85 {
		t.Fatalf("root.Left.Value expected to have 4.85, got %f", root.Left.Value)
	}
	if root.Left.Left != nil {
		t.Fatalf("root.Left.Left expected nil, got %v", root.Left.Left)
	}
	if root.Right.Value != 5.77 {
		t.Fatalf("root.Right.Value expected to have 5.77, got %f", root.Right.Value)
	}
	if root.Right.Right.Value != 7.999 {
		t.Fatalf("root.Right.Right.Value expected to have 7.999, got %f", root.Right.Right.Value)
	}

	s1 := search(root, 5.77)
	if s1 == nil {
		t.Fatalf("5.77 not found, got %v", s1)
	}
	if s1.Value != 5.77 {
		t.Fatalf("s1 expected to have 5.77, got %f", s1.Value)
	}

	s2 := searchClosest(root, 7.7)
	if s2 == nil {
		t.Fatalf("closest value to 7.7 not found, got %v", s2)
	}
	if s2.Value != 7.999 {
		t.Fatalf("s2 expected to have 7.999, got %f", s2.Value)
	}

	s3 := searchClosest(root, 777.7)
	if s3 == nil {
		t.Fatalf("closest value to 777.7 not found, got %v", s3)
	}
	if s3.Value != 7.999 {
		t.Fatalf("s3 expected to have 7.999, got %f", s3.Value)
	}
	if s3.Idx != 3 {
		t.Fatalf("s3.Idx expected to have 3, got %f", s3.Idx)
	}

	s4 := searchClosest(root, 5.66)
	if s4 == nil {
		t.Fatalf("closest value to 5.66 not found, got %v", s4)
	}
	if s4.Value != 5.77 {
		t.Fatalf("s4 expected to have 5.77, got %f", s4.Value)
	}

	s5 := searchClosest(root, 5.21)
	if s5 == nil {
		t.Fatalf("closest value to 5.21 not found, got %v", s5)
	}
	if s5.Value != 5.22 {
		t.Fatalf("s5 expected to have 5.22, got %f", s5.Value)
	}

	s6 := searchClosest(root, 5.22)
	if s6 == nil {
		t.Fatalf("closest value to 5.22 not found, got %v", s6)
	}
	if s6.Value != 5.22 {
		t.Fatalf("s6 expected to have 5.22, got %f", s6.Value)
	}

	s7 := searchClosest(root, 4.85)
	if s7 == nil {
		t.Fatalf("closest value to 4.85 not found, got %v", s7)
	}
	if s7.Value != 4.85 {
		t.Fatalf("s7 expected to have 4.85, got %f", s7.Value)
	}
}

func Test_boundaries(t *testing.T) {
	all := []int64{15, 5, 1}
	bds := buildBoundaries(all)

	bd1 := bds.findBoundary(2)
	exp1 := boundary{lower: 1, lowerIdx: 2, upper: 5, upperIdx: 1}
	if !reflect.DeepEqual(bd1, exp1) {
		t.Fatalf("boundary expected %+v, got %+v", exp1, bd1)
	}

	bd2 := bds.findBoundary(100)
	exp2 := boundary{lower: 15, lowerIdx: 0, upper: 0, upperIdx: -1}
	if !reflect.DeepEqual(bd2, exp2) {
		t.Fatalf("boundary expected %+v, got %+v", exp2, bd2)
	}

	bd3 := bds.findBoundary(0)
	exp3 := boundary{lower: 0, lowerIdx: -1, upper: 1, upperIdx: 2}
	if !reflect.DeepEqual(bd3, exp3) {
		t.Fatalf("boundary expected %+v, got %+v", exp3, bd3)
	}

	bd4 := bds.findBoundary(13)
	exp4 := boundary{lower: 5, lowerIdx: 1, upper: 15, upperIdx: 0}
	if !reflect.DeepEqual(bd4, exp4) {
		t.Fatalf("boundary expected %+v, got %+v", exp4, bd4)
	}

	bds.add(12)

	bd4 = bds.findBoundary(13)
	exp4 = boundary{lower: 12, lowerIdx: 3, upper: 15, upperIdx: 0}
	if !reflect.DeepEqual(bd4, exp4) {
		t.Fatalf("boundary expected %+v, got %+v", exp4, bd4)
	}

	bd5 := bds.findBoundary(12)
	exp5 := boundary{lower: 12, lowerIdx: 3, upper: 12, upperIdx: 3}
	if !reflect.DeepEqual(bd5, exp5) {
		t.Fatalf("boundary expected %+v, got %+v", exp5, bd5)
	}
}

func TestTree(t *testing.T) {
	nums := []float64{5.22, 4.85, 5.77, 7.999}
	tr := NewBinaryTree(nums)
	idx, v := tr.Closest(4.85)
	if idx != 1 || v != 4.85 {
		t.Fatalf("idx, v expected 1, 4.85 / got %d, %f", idx, v)
	}
	idx, v = tr.Closest(5.66)
	if idx != 2 || v != 5.77 {
		t.Fatalf("idx, v expected 2, 5.77 / got %d, %f", idx, v)
	}
	idx, v = tr.Closest(7.9999)
	if idx != 3 || v != 7.999 {
		t.Fatalf("idx, v expected 3, 7.999 / got %d, %f", idx, v)
	}
}
