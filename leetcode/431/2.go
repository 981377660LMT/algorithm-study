package main

func calculateScore(s string) int64 {
	var res int64

	stack := make([][]int, 26)
	for i := 0; i < len(s); i++ {
		curIndex := s[i] - 'a'
		mIndex := 25 - curIndex
		if len(stack[mIndex]) > 0 {
			last := len(stack[mIndex]) - 1
			j := stack[mIndex][last]
			stack[mIndex] = stack[mIndex][:last]
			res += int64(i - j)
		} else {
			stack[curIndex] = append(stack[curIndex], i)
		}
	}
	return res
}
