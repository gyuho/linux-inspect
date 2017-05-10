package inspect

import (
	"math"
	"sort"
)

// BinarySearchInt64 binary-searches the int64 slice
// and returns the index of the matching element.
// So input slice must be sorted.
// It returns -1 if not found.
func BinarySearchInt64(nums []int64, v int64) int {
	lo := 0
	hi := len(nums) - 1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if nums[mid] < v {
			lo = mid + 1 // keep searching on right-subtree
			continue
		}

		if nums[mid] > v {
			hi = mid - 1 // keep searching on left-subtree
			continue
		}

		return mid
	}
	return -1
}

// Tree defines binary search tree.
type Tree interface {
	Closest(v float64) (index int, value float64)
}

// NewBinaryTree builds a new binary search tree.
// The original slice won't be sorted.
func NewBinaryTree(nums []float64) Tree {
	if len(nums) == 0 {
		return nil
	}

	root := newFloat64Node(0, nums[0])
	for i := range nums {
		if i == 0 {
			continue
		}
		insert(root, i, nums[i])
	}
	return root
}

// NewBinaryTreeInt64 builds a new binary search tree.
// The original slice won't be sorted.
func NewBinaryTreeInt64(nums []int64) Tree {
	fs := make([]float64, len(nums))
	for i := range nums {
		fs[i] = float64(nums[i])
	}
	return NewBinaryTree(fs)
}

func (root *float64Node) Closest(v float64) (index int, value float64) {
	nd := searchClosest(root, v)
	return nd.Idx, nd.Value
}

// float64Node represents binary search tree
// to find the closest float64 value.
type float64Node struct {
	Idx   int
	Value float64
	Left  *float64Node
	Right *float64Node
}

// newFloat64Node returns a new float64Node.
func newFloat64Node(idx int, v float64) *float64Node {
	return &float64Node{Idx: idx, Value: v}
}

// insert inserts a value to the binary search tree.
// For now, it assumes that values are unique.
func insert(root *float64Node, idx int, v float64) *float64Node {
	if root == nil {
		return newFloat64Node(idx, v)
	}

	if root.Value > v {
		root.Left = insert(root.Left, idx, v)
	} else {
		root.Right = insert(root.Right, idx, v)
	}

	return root
}

// search searches a value in the binary search tree.
func search(root *float64Node, v float64) *float64Node {
	if root == nil {
		return nil
	}

	if root.Value == v {
		return root
	}

	if root.Value > v {
		return search(root.Left, v)
	}

	return search(root.Right, v)
}

// searchClosest searches the closest value in the binary search tree.
func searchClosest(root *float64Node, v float64) *float64Node {
	if root == nil {
		return nil
	}

	var child *float64Node
	if root.Value > v {
		child = searchClosest(root.Left, v)
	} else {
		child = searchClosest(root.Right, v)
	}

	// no children, just return root
	if child == nil {
		return root
	}

	rootDiff := math.Abs(float64(root.Value - v))
	childDiff := math.Abs(float64(child.Value - v))
	if rootDiff < childDiff {
		// diff with root is smaller
		return root
	}

	return child
}

// boundary is the pair of values in a boundary.
type boundary struct {
	// index of 'lower' in the original slice
	lower    int64
	lowerIdx int

	// index of 'upper' in the original slice
	upper    int64
	upperIdx int
}

type boundaries struct {
	// store original slice as well
	// to return the index
	numsOrig    []int64
	num2OrigIdx map[int64]int

	numsSorted    []int64
	num2SortedIdx map[int64]int

	tr Tree
}

func buildBoundaries(nums []int64) *boundaries {
	num2OrigIdx := make(map[int64]int)
	for i := range nums {
		num2OrigIdx[nums[i]] = i
	}
	numsOrig := make([]int64, len(nums))
	copy(numsOrig, nums)

	tr := NewBinaryTreeInt64(nums)

	sort.Sort(int64Slice(nums))
	num2SortedIdx := make(map[int64]int)
	for i := range nums {
		num2SortedIdx[nums[i]] = i
	}

	return &boundaries{
		numsOrig:      numsOrig,
		num2OrigIdx:   num2OrigIdx,
		numsSorted:    nums,
		num2SortedIdx: num2SortedIdx,
		tr:            tr,
	}
}

// adds a second to boundaries
// and rebuild the binary tree
func (bf *boundaries) add(sec int64) {
	bf.numsOrig = append(bf.numsOrig, sec)
	bf.num2OrigIdx[sec] = len(bf.numsOrig)

	bf.numsSorted = append(bf.numsSorted, sec)

	// re-sort
	bf.tr = NewBinaryTreeInt64(bf.numsSorted)
	sort.Sort(int64Slice(bf.numsSorted))

	num2SortedIdx := make(map[int64]int)
	for i := range bf.numsSorted {
		num2SortedIdx[bf.numsSorted[i]] = i
	}
	bf.num2SortedIdx = num2SortedIdx
}

// returns the boundary with closest upper, lower value.
// returns the index of the value if found.
func (bf *boundaries) findBoundary(missingSecond int64) (bd boundary) {
	idxOrig, vOrig := bf.tr.Closest(float64(missingSecond))
	valOrig := int64(vOrig)
	if valOrig == missingSecond {
		bd.lower = valOrig
		bd.lowerIdx = idxOrig
		bd.upper = valOrig
		bd.upperIdx = idxOrig
		return
	}

	// use the idx in sorted!
	idxx := bf.num2SortedIdx[valOrig]

	if missingSecond > valOrig {
		bd.lower = valOrig
		bd.lowerIdx = idxOrig

		// valOrig is the lower bound, we need to find another upper value
		// continue search in right-half
		// (assume 'nums' is sorted)
		for j := idxx + 1; j < len(bf.numsSorted); j++ {
			if bf.numsSorted[j] > missingSecond {
				// found upper bound
				bd.upper = bf.numsSorted[j]
				bd.upperIdx = bf.num2OrigIdx[bf.numsSorted[j]]
				return
			}
		}

		bd.upper = 0
		bd.upperIdx = -1
		return
	}
	bd.upper = valOrig
	bd.upperIdx = idxOrig

	// valOrig is the upper bound, we need to find another lower value
	// continue search in left-half
	// (assume 'nums' is sorted)
	for j := idxx - 1; j >= 0; j-- {
		if bf.numsSorted[j] < missingSecond {
			// found lower bound
			bd.lower = bf.numsSorted[j]
			bd.lowerIdx = bf.num2OrigIdx[bf.numsSorted[j]]
			return
		}
	}

	bd.lower = 0
	bd.lowerIdx = -1
	return
}

type int64Slice []int64

func (s int64Slice) Len() int           { return len(s) }
func (s int64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s int64Slice) Less(i, j int) bool { return s[i] < s[j] }
