# # For how many positive integers at most N is the product of the digits at most K
# # !1-n中 各位上的积不超过k的数的个数
# n<=1e9 k<=1e18

from functools import lru_cache


def cal(upper: int, k: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, mul: int) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界"""
        if pos == n:
            return not hasLeadingZero and mul <= k

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), 1)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), mul * cur)
        return res

    nums = list(map(int, str(upper)))
    n = len(nums)
    return dfs(0, True, True, 1)


class Solution:
    def main(self, n: int, k: int) -> int:
        """不超过n的正整数中,有多少个数满足`各位之积不超过k`"""
        return cal(n, k) - cal(0, k)


print(Solution().main(n=13, k=2))  # 5
print(Solution().main(n=100, k=80))  # 99
print(Solution().main(n=1000000000000000000, k=1000000000))  # 841103275147365677


# 为什么能用乘积作为状态?
# 考虑质因子 2 3 5 7 的个数
# !2的个数不超过60 3的个数不超过38 5的个数不超过26 7的个数不超过22 状态数少了很多
