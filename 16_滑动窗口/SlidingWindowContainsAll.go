package main

const INF int = 2e18

// 面试题 17.18. 最短超串
// https://leetcode.cn/problems/shortest-supersequence-lcci/description/
// 找到长数组中包含短数组所有的元素的最短子数组，其出现顺序无关紧要。
// 返回最短子数组的左端点和右端点，如有多个满足条件的子数组，返回左端点最小的一个。若不存在，返回空数组。
func shortestSeq(big []int, small []int) []int {
	if len(big) < len(small) {
		return nil
	}

	S := NewSlidingWindowContainsAll(int32(len(small)), func(i int32) int { return small[i] })
	resLen, resLeft, left, n := INF, -1, 0, len(big)
	for right := 0; right < n; right++ {
		S.Add(big[right])
		for left <= right && S.ContainsAll() {
			if len := right - left + 1; len < resLen || (len == resLen && left < resLeft) {
				resLen, resLeft = len, left
			}

			S.Discard(big[left])
			left++
		}
	}

	if resLen == INF {
		return nil
	}
	return []int{resLeft, resLeft + resLen - 1}
}

// 3298. 统计重新排列后包含另一个字符串的子字符串数目 II
// https://leetcode.cn/problems/count-substrings-that-can-be-rearranged-to-contain-a-string-ii/description/
func validSubstringCount(word1 string, word2 string) int64 {
	S := NewSlidingWindowContainsAll(int32(len(word2)), func(i int32) byte { return word2[i] })
	res, left, n := 0, 0, len(word1)
	for right := 0; right < n; right++ {
		S.Add(word1[right])
		for left <= right && S.ContainsAll() {
			S.Discard(word1[left])
			left++
		}
		res += left
	}
	return int64(res)
}

// 支持新增元素、删除元素、快速判断容器内是否包含了指定的所有元素.
type SlidingWindowContainsAll[T comparable] struct {
	missingCount int32
	missing      map[T]int32
}

func NewSlidingWindowContainsAll[T comparable](n int32, supplier func(i int32) T) *SlidingWindowContainsAll[T] {
	missing := make(map[T]int32)
	for i := int32(0); i < n; i++ {
		missing[supplier(i)]++
	}
	missingCount := int32(len(missing))
	return &SlidingWindowContainsAll[T]{missingCount: missingCount, missing: missing}
}

func (s *SlidingWindowContainsAll[T]) Add(v T) bool {
	if pre, has := s.missing[v]; !has {
		return false
	} else {
		s.missing[v] = pre - 1
		if pre == 1 {
			s.missingCount--
		}
		return true
	}
}

func (s *SlidingWindowContainsAll[T]) Discard(v T) bool {
	if pre, has := s.missing[v]; !has {
		return false
	} else {
		s.missing[v] = pre + 1
		if pre == 0 {
			s.missingCount++
		}
		return true
	}
}

func (s *SlidingWindowContainsAll[T]) ContainsAll() bool {
	return s.missingCount == 0
}
