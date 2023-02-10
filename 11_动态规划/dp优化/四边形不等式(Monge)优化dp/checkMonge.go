package main

// 检查[0,n]内f是否具有monge性质
func checkMonge(n int, f func(i, j int) int) bool {
	for l := 0; l <= n; l++ {
		for k := 0; k < l; k++ {
			for j := 0; j < k; j++ {
				for i := 0; i < j; i++ {
					lhs := f(i, l) + f(j, k)
					rhs := f(i, k) + f(j, l)
					if lhs < rhs {
						return false
					}
				}
			}
		}
	}
	return true
}
