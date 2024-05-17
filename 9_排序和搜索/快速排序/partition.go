package main

import "fmt"

// partition_func does one quicksort partition.
// Let p = data[pivot]
// Moves elements in data[a:b] around, so that data[i]<p and data[j]>=p for i<newpivot and j>newpivot.
// On return, data[newpivot] = p
func partition_func(less func(i, j int) bool, swap func(i, j int), a, b, pivot int) (newpivot int, alreadyPartitioned bool) {
	swap(a, pivot)
	i, j := a+1, b-1 // i and j are inclusive of the elements remaining to be partitioned

	for i <= j && less(i, a) {
		i++
	}
	for i <= j && !less(j, a) {
		j--
	}
	if i > j {
		swap(j, a)
		return j, true
	}
	swap(i, j)
	i++
	j--

	for {
		for i <= j && less(i, a) {
			i++
		}
		for i <= j && !less(j, a) {
			j--
		}
		if i > j {
			break
		}
		swap(i, j)
		i++
		j--
	}
	swap(j, a)
	return j, false
}

func main() {
	nums := []int{3, 2, 1, 5, 6, 4}

	newPivot, sorted := partition_func(
		func(i, j int) bool { return nums[i] < nums[j] },
		func(i, j int) { nums[i], nums[j] = nums[j], nums[i] },
		0, len(nums), 3,
	)

	fmt.Println(nums, newPivot, sorted)
}
