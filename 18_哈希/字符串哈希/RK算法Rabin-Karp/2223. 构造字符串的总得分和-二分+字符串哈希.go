package main

func sumScores(s string) int64 {
	n := len(s)
	hasher := BiStringHasher(s, 1e8+7, 1e9+7, 131, 13131, 0, 0)
	countPre := func(curLen, start int) int {
		left, right := 1, curLen
		for left <= right {
			mid := (left + right) >> 1
			hash00, hash01 := hasher(start, start+mid)
			hash10, hash11 := hasher(0, mid)
			if hash00 == hash10 && hash01 == hash11 {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		return right
	}

	res := 0
	for i := 1; i < n+1; i++ {
		count := countPre(i, n-i)
		res += count
	}

	return int64(res)

}

func StringHasher(s string, mod int, base int, offset int) func(left int, right int) int {
	prePow := make([]int, len(s)+1)
	prePow[0] = 1
	preHash := make([]int, len(s)+1)
	for i, v := range s {
		prePow[i+1] = (prePow[i] * base) % mod
		preHash[i+1] = (preHash[i]*base + int(v) - offset) % mod
	}

	sliceHash := func(left, right int) int {
		if left >= right {
			return 0
		}
		return (preHash[right] - preHash[left]*prePow[right-left]%mod + mod) % mod
	}

	return sliceHash
}

func BiStringHasher(s string, mod1, mod2, base1, base2, offset1, offset2 int) func(left int, right int) (hash1, hash2 int) {
	hasher1 := StringHasher(s, mod1, base1, offset1)
	hasher2 := StringHasher(s, mod2, base2, offset2)

	sliceHash := func(left, right int) (hash1, hash2 int) {
		if left >= right {
			return 0, 0
		}

		hash1 = hasher1(left, right)
		hash2 = hasher2(left, right)
		return
	}

	return sliceHash
}
