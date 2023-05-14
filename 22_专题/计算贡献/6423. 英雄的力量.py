# 6423. 英雄的力量
# !求所有非空子集(最小值*最大值的平方)的和
# !排序，枚举最大值


from typing import List

MOD = int(1e9 + 7)


class Solution:
    def sumOfPower(self, nums: List[int]) -> int:
        nums.sort()
        res = 0
        curSum = 0
        for num in nums:
            res += (num + curSum) * num * num
            res %= MOD
            curSum = (curSum * 2 + num) % MOD
        return res


assert Solution().sumOfPower(nums=[2, 1, 4]) == 141
