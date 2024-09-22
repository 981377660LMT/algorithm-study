package main

// 给你两个字符串 word1 和 word2 。

// 如果一个字符串 x 重新排列后，word2 是重排字符串的 前缀 ，那么我们称字符串 x 是 合法的 。

// 请你返回 word1 中 合法 子字符串 的数目。
func validSubstringCount(word1 string, word2 string) int64 {
	counter := [26]int32{}
	for _, c := range word2 {
		counter[c-'a']++
	}

	isValid := func() bool {
		for _, c := range counter {
			if c > 0 {
				return false
			}
		}
		return true
	}

	res, left, n := 0, 0, len(word1)
	for right := 0; right < n; right++ {
		counter[word1[right]-'a']--
		for left < n && isValid() {
			counter[word1[left]-'a']++
			left++
		}
		res += left
	}
	return int64(res)
}
