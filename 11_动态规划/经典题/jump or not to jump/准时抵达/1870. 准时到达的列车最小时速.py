from typing import List
from math import ceil

# 返回能满足你准时到达办公室所要求全部列车的 最小正整数 时速（单位：千米每小时），如果无法准时到达，则返回 -1 。


class Solution:
    def minSpeedOnTime(self, dist: List[int], hour: float) -> int:
        def check(mid):
            res = 0
            for i in range(len(dist) - 1):
                res += ceil(dist[i] / mid)
            res += dist[-1] / mid
            return res <= hour

        if len(dist) - 1 >= hour:
            return -1

        left, right = 1, int(1e20)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1

        return left


print(Solution().minSpeedOnTime(dist=[1, 3, 2], hour=6))
# 输出：1
# 解释：速度为 1 时：
# - 第 1 趟列车运行需要 1/1 = 1 小时。
# - 由于是在整数时间到达，可以立即换乘在第 1 小时发车的列车。第 2 趟列车运行需要 3/1 = 3 小时。
# - 由于是在整数时间到达，可以立即换乘在第 4 小时发车的列车。第 3 趟列车运行需要 2/1 = 2 小时。
# - 你将会恰好在第 6 小时到达。

