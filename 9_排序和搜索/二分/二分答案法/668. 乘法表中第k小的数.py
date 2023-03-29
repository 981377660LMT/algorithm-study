# m,n<=3e4
# 668. 乘法表中第k小的数


class Solution:
    def findKthNumber(self, m: int, n: int, k: int) -> int:
        """时间复杂度O(min(m,n)logmn)"""
        if m > n:  # 优化
            m, n = n, m

        def countNGT(mid: int) -> int:
            """有多少个不超过mid的候选"""
            res = 0
            for r in range(1, m + 1):
                res += min(mid // r, n)
            return res

        left, right = 0, int(1e18)
        while left <= right:
            mid = (left + right) // 2
            if countNGT(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left
