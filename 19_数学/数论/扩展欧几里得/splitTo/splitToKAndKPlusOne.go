package main

// SplitToKAndKPlusOne 将 num 拆分成 k 和 k+1 的和，使得拆分的个数最(多/少).
//  minimize 是否使得拆分的个数最少. 默认为最少(true).
//  count1和count2分别是拆分成k和k+1的个数，ok表示是否可以拆分.
func SplitToKAndKPlusOne(num, k int, minimize bool) (count1, count2 int, ok bool) {
	if minimize {
		count2 = (num + k) / (k + 1)
		diff := (k+1)*count2 - num
		if diff > count2 {
			return 0, 0, false
		}
		return diff, count2 - diff, true
	}

	count1 = num / k
	diff := num - k*count1
	if diff > count1 {
		return 0, 0, false
	}
	return count1 - diff, diff, true
}
