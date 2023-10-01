// type DiffArray struct {
// 	diff  []int
// 	dirty bool
// }

// func NewDiffArray(n int) *DiffArray {
// 	return &DiffArray{
// 		diff: make([]int, n+1),
// 	}
// }

// func (d *DiffArray) Add(start, end, delta int) {
// 	if start < 0 {
// 		start = 0
// 	}
// 	if end >= len(d.diff) {
// 		end = len(d.diff) - 1
// 	}
// 	if start >= end {
// 		return
// 	}
// 	d.dirty = true
// 	d.diff[start] += delta
// 	d.diff[end] -= delta
// }

// func (d *DiffArray) Build() {
// 	if d.dirty {
// 		preSum := make([]int, len(d.diff))
// 		for i := 1; i < len(d.diff); i++ {
// 			preSum[i] = preSum[i-1] + d.diff[i]
// 		}
// 		d.diff = preSum
// 		d.dirty = false
// 	}
// }

// func (d *DiffArray) Get(pos int) int {
// 	d.Build()
// 	return d.diff[pos]
// }

// func (d *DiffArray) GetAll() []int {
// 	d.Build()
// 	return d.diff[:len(d.diff)-1]
// }

// type DiffMap struct {
// 	diff       map[int]int
// 	sortedKeys []int
// 	preSum     []int
// 	dirty      bool
// }

// func NewDiffMap() *DiffMap {
// 	return &DiffMap{
// 		diff: make(map[int]int),
// 	}
// }

// func (d *DiffMap) Add(start, end, delta int) {
// 	if start >= end {
// 		return
// 	}
// 	d.dirty = true
// 	d.diff[start] += delta
// 	d.diff[end] -= delta
// }

// func (d *DiffMap) Build() {
// 	if d.dirty {
// 		d.sortedKeys = make([]int, 0, len(d.diff))
// 		for key := range d.diff {
// 			d.sortedKeys = append(d.sortedKeys, key)
// 		}
// 		sort.Ints(d.sortedKeys)
// 		d.preSum = make([]int, len(d.sortedKeys)+1)
// 		for i, key := range d.sortedKeys {
// 			d.preSum[i+1] = d.preSum[i] + d.diff[key]
// 		}
// 		d.dirty = false
// 	}
// }

// func (d *DiffMap) Get(pos int) int {
// 	d.Build()
// 	return d.preSum[sort.SearchInts(d.sortedKeys, pos+1)]
// }

class DiffArray {}

class DiffMap {}

export { DiffArray, DiffMap }
