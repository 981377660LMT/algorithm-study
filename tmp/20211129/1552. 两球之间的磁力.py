from typing import List

# 想把 m 个球放到这些篮子里，使得任意两球间 最小磁力 最大。
# 已知两个球如果分别位于 x 和 y ，那么它们之间的磁力为 |x - y| 。
# 请你返回最大化的最小磁力。

# 2 <= n <= 10^5 暗示nlogn
# 最左能力二分
class Solution:
    def maxDistance(self, position: List[int], m: int) -> int:
        n = len(position)
        position.sort()

        # 最小磁力不可超过mid，需要多少个球
        def count(mid: int) -> int:
            res, cur = 1, position[0]
            for i in range(1, n):
                if position[i] - cur > mid:
                    res += 1
                    cur = position[i]
            return res

        l, r = 0, position[-1] - position[0]
        while l <= r:
            mid = (l + r) >> 1
            if count(mid) >= m:
                l = mid + 1
            else:
                r = mid - 1

        return l


print(Solution().maxDistance(position=[1, 2, 3, 4, 7], m=3))
# 输出：3
# 解释：将 3 个球分别放入位于 1，4 和 7 的三个篮子，两球间的磁力分别为 [3, 3, 6]。最小磁力为 3 。我们没办法让最小磁力大于 3 。

