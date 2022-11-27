# 1 <= lowLimit <= highLimit <= 1e5
# !求各位数字之和出现次数的最大值


from functools import lru_cache


def cal(upper: int, target: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: int, isLimit: bool, curSum: int) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界"""
        if curSum > target:  # !剪枝
            return 0
        if pos == len(nums):
            return 1 if curSum == target else 0

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), curSum)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), curSum + cur)
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, True, True, 0)


class Solution:
    def countBalls(self, lowLimit: int, highLimit: int) -> int:
        return max(
            cal(highLimit, target) - cal(lowLimit - 1, target)
            for target in range(1, 9 * len(str(highLimit)) + 1)
        )
