from typing import List

# 想把 m 个球放到这些篮子里，使得任意两球间 最小磁力 最大。
# 已知两个球如果分别位于 x 和 y ，那么它们之间的磁力为 |x - y| 。
# 请你返回最大化的最小磁力。

# 2 <= n <= 10^5 暗示nlogn
# 最左能力二分


class Solution:
    def maxDistance(self, position: List[int], m: int) -> int:
        def check(mid: int) -> int:
            """最小磁力为mid时是否能放m个球"""
            res, pre = 1, position[0]
            for i in range(1, n):
                if position[i] - pre >= mid:
                    res += 1
                    pre = position[i]
            return res >= m

        n = len(position)
        position.sort()

        left, right = 1, int(1e10)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        return right


print(Solution().maxDistance(position=[1, 2, 3, 4, 7], m=3))
# 输出：3
# 解释：将 3 个球分别放入位于 1，4 和 7 的三个篮子，两球间的磁力分别为 [3, 3, 6]。最小磁力为 3 。我们没办法让最小磁力大于 3 。

