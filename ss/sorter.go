package ss

import "sort"

// by returns a multiSorter that sorts using the less functions
func by(rows [][]string, lesses ...lessFunc) *multiSorter {
	return &multiSorter{
		data: rows,
		less: lesses,
	}
}

// lessFunc compares between two string slices.
type lessFunc func(p1, p2 *[]string) bool

func makeAscendingFunc(idx int) func(row1, row2 *[]string) bool {
	return func(row1, row2 *[]string) bool {
		return (*row1)[idx] < (*row2)[idx]
	}
}

// multiSorter implements the Sort interface,
// sorting the two dimensional string slices within.
type multiSorter struct {
	data [][]string
	less []lessFunc
}

// Sort sorts the rows according to lessFunc.
func (ms *multiSorter) Sort(rows [][]string) {
	sort.Sort(ms)
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.data)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.data[i], ms.data[j] = ms.data[j], ms.data[i]
}

// Less is part of sort.Interface.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.data[i], &ms.data[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q
			return true
		case less(q, p):
			// p > q
			return false
		}
		// p == q; try next comparison
	}
	return ms.less[k](p, q)
}
