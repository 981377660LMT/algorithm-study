from typing import List
from itertools import groupby

# 2 <= seats.length <= 2 * 104

# 其实就是计算最大连续0的长度，头和尾特殊计算
# 1.计算座位到最近的人的最大距离 两个辅助数组
# 2.按零分组


class Solution:
    def maxDistToClosest(self, seats: List[int]) -> int:
        n = len(seats)
        leftDist, rightDist = [n] * n, [n] * n
        for i in range(n):
            if seats[i] == 1:
                leftDist[i] = 0
            elif i > 0:
                leftDist[i] = leftDist[i - 1] + 1

        for i in range(n - 1, -1, -1):
            if seats[i] == 1:
                rightDist[i] = 0
            elif i < n - 1:
                rightDist[i] = rightDist[i + 1] + 1

        return max(min(leftDist[i], rightDist[i]) for i, seat in enumerate(seats) if seat == 0)

    # 如果两人之间有连续 K 个空座位，那么其中存在至少一个座位到两边最近的人的距离为 (K+1) / 2。
    def maxDistToClosest2(self, seats: List[int]) -> int:
        res = 0
        for zeroOrOne, group in groupby(seats):
            if zeroOrOne == 0:
                k = len(list(group))
                res = max(res, (k + 1) >> 1)

        return max(res, seats.index(1), seats[::-1].index(1))


print(Solution().maxDistToClosest(seats=[1, 0, 0, 0, 1, 0, 1]))
# 输出：2
# 解释：
# 如果亚历克斯坐在第二个空位（seats[2]）上，他到离他最近的人的距离为 2 。
# 如果亚历克斯坐在其它任何一个空位上，他到离他最近的人的距离为 1 。
# 因此，他到离他最近的人的最大距离是 2 。
