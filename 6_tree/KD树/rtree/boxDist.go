// BoxDist 计算两个矩形之间的最小距离。如果两个矩形相交或相邻，则距离为零；
// 否则，距离是两个矩形之间的欧几里得距离（通常是矩形边缘之间的最短距离）。

package main

type numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type rect[N numeric] struct {
	min [2]N
	max [2]N
}

// 高阶函数，用于计算 R-树中条目与目标矩形之间的距离.
// 用户可以通过提供自定义的 itemDist 函数，定义自己的距离计算方式.
func BoxDist[N numeric, T any](targetMin, targetMax [2]N, itemDist func(min, max [2]N, data T) N) (dist func(min, max [2]N, data T, item bool) N) {
	targ := rect[N]{targetMin, targetMax}

	return func(min, max [2]N, data T, item bool) (dist N) {
		if item && itemDist != nil {
			return itemDist(min, max, data)
		}
		return targ.boxDist(&rect[N]{min, max})
	}
}

// 矩形距离.
func (r *rect[N]) boxDist(b *rect[N]) N {
	var dist N
	squared := fmax(r.min[0], b.min[0]) - fmin(r.max[0], b.max[0])
	if squared > 0 {
		dist += squared * squared
	}
	squared = fmax(r.min[1], b.min[1]) - fmin(r.max[1], b.max[1])
	if squared > 0 {
		dist += squared * squared
	}
	return dist
}

func fmin[N numeric](a, b N) N {
	if a < b {
		return a
	}
	return b
}

func fmax[N numeric](a, b N) N {
	if a > b {
		return a
	}
	return b
}
