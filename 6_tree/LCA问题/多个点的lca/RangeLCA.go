package main

// 区间LCA.
//
//	points 顶点数组.
//	getLCA LCA实现.
//	返回一个查询 points[start:end) lca 的函数.
func RangeLCA(
	points []int,
	getLCA func(u, v int) int,
) func(start, end int) int {
	n := 1
	for n < len(points) {
		n <<= 1
	}
	seg := make([]int32, n<<1)
	for i := 0; i < len(points); i++ {
		seg[n+i] = int32(points[i])
	}
	for i := n - 1; i >= 0; i-- {
		seg[i] = int32(getLCA(int(seg[i<<1]), int(seg[(i<<1)|1])))
	}
	lca := func(u, v int32) int32 {
		if u == -1 || v == -1 {
			if u == -1 {
				return v
			}
			return u
		}
		return int32(getLCA(int(u), int(v)))
	}
	query := func(start, end int) int {
		res := int32(-1)
		for ; start > 0 && start+(start&-start) <= end; start += start & -start {
			res = lca(res, seg[(n+start)/(start&-start)])
		}
		for ; start < end; end -= end & -end {
			res = lca(res, seg[(n+end)/(end&-end)-1])
		}
		return int(res)
	}
	return query
}
