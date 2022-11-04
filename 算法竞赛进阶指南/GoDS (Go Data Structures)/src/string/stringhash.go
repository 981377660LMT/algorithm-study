package cmnx

func main() {

}
hash := func(s string) {
	// 注意：由于哈希很容易被卡，能用其它方法实现尽量用其它方法
	const prime uint64 = 1e8 + 7
	powP := make([]uint64, len(s)+1) // powP[i] = prime^i
	powP[0] = 1
	preHash := make([]uint64, len(s)+1) // preHash[i] = hash(s[:i]) 前缀哈希
	for i, b := range s {
		powP[i+1] = powP[i] * prime
		preHash[i+1] = preHash[i]*prime + uint64(b) // 本质是秦九韶算法
	}

	// 计算子串 s[l:r] 的哈希   0<=l<=r<=len(s)
	// 空串的哈希值为 0
	subHash := func(l, r int) uint64 { return preHash[r] - preHash[l]*powP[r-l] }
	_ = subHash
}


// !需要双哈希