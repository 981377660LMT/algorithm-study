# 3279. 活塞占据的最大总区域
# https://leetcode.cn/problems/maximum-total-area-occupied-by-pistons/solutions/2906015/python3chai-fen-shu-zu-by-arnold-sb6ffyl-2ppp/
# 一台旧车的引擎中有一些活塞，我们想要计算活塞 下方 的 最大 区域。
#
# 给定：
#
# 一个整数 height，表示活塞 最大 可到达的高度。
# 一个整数数组 positions，其中 positions[i] 是活塞 i 的当前位置，等于其 下方 的当前区域。
# 一个字符串 directions，其中 directions[i] 是活塞 i 的当前移动方向，'U' 表示向上，'D' 表示向下。
# 每一秒：
# 每个活塞向它的当前方向移动 1 单位。即如果方向向上，positions[i] 增加 1。
# 如果一个活塞到达了其中一个终点，即 positions[i] == 0 或 positions[i] == height，它的方向将会改变。
# 返回所有活塞下方的最大可能区域。
#
# 每个活塞都会在 2 * height 秒后会回到初始点，
# 这个过程中会有两个时间点变换方向。
# 每秒面积变化为所有活塞朝上的个数减去朝下的个数。
# 记录下每个活塞变换方向的时间点，按时间顺序遍历这些时间点，模拟即可。


from collections import defaultdict
from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxArea(self, height: int, positions: List[int], directions: str) -> int:
        dirs = list(directions)
        diff = defaultdict(int)  # 维护向上的活塞个数减去向下的活塞个数

        def add(left: int, right: int, delta: int):
            diff[left] += delta
            diff[right] -= delta

        curSum, curDiff = 0, 0
        for i, pos in enumerate(positions):
            curSum += pos
            if pos == 0:
                dirs[i] = "U"
            if pos == height:
                dirs[i] = "D"

            if dirs[i] == "U":
                curDiff += 1
                add(height - pos, 2 * height - pos, -2)
            else:
                curDiff -= 1
                add(pos, height + pos, 2)

        res = curSum
        preT = 0
        for t in sorted(diff):
            curSum += (t - preT) * curDiff
            curDiff += diff[t]
            preT = t
            res = max2(res, curSum)
        return res
