package sortedlist

func containsNearbyAlmostDuplicate(nums []int, indexDiff int, valueDiff int) bool {
	sl := NewSortedList(func(a, b interface{}) int {
		return a.(int) - b.(int)
	}, 16)

	for right := 0; right < len(nums); right++ {
		if right-indexDiff-1 >= 0 {
			sl.Discard(nums[right-indexDiff-1])
		}

		pos1 := sl.BisectRight(nums[right]-valueDiff) - 1
		if pos1 >= 0 && abs(nums[right]-sl.At(pos1).(int)) <= valueDiff {
			return true
		}

		pos2 := sl.BisectRight(nums[right]+valueDiff) - 1
		if pos2 >= 0 && abs(nums[right]-sl.At(pos2).(int)) <= valueDiff {
			return true
		}

		sl.Add(nums[right])
	}

	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
