// 3296. 移山所需的最少秒数-二分套二分
// https://leetcode.cn/problems/minimum-number-of-seconds-to-make-mountain-height-zero/description/
// 给你一个整数 mountainHeight 表示山的高度。
// 同时给你一个整数数组 workerTimes，表示工人们的工作时间（单位：秒）。
// 工人们需要 同时 进行工作以 降低 山的高度。对于工人 i :
// 山的高度降低 x，需要花费 workerTimes[i] + workerTimes[i] * 2 + ... + workerTimes[i] * x 秒。例如：
// 山的高度降低 1，需要 workerTimes[i] 秒。
// 山的高度降低 2，需要 workerTimes[i] + workerTimes[i] * 2 秒，依此类推。
// 返回一个整数，表示工人们使山的高度降低到 0 所需的 最少 秒数。

package main

const INF int = 1e16

func minNumberOfSeconds(mountainHeight int, workerTimes []int) int64 {
	// 工作x秒，降低高度之和能否>=mountainHeight.
	// 对每个人，找到降低的最大高度，使得工作时间<=x.
	check := func(x int) bool {
		res := 0
		for _, t := range workerTimes {
			right := MaxRight(
				0,
				func(right int) bool {
					right--
					return t*(right+1)*right/2 <= x
				},
				mountainHeight+1,
			)
			res += right - 1
		}

		return res >= mountainHeight
	}

	res := MinLeft(INF, check, 0)
	return int64(res)
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含,使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}
