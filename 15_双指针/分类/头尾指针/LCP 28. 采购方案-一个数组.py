# 假定小力仅购买两个零件，要求购买零件的花费不超过预算，请问他有多少种采购方案
# 注意(a,b)和(b,a)是同一种方案


from typing import List


MOD = int(1e9 + 7)


class Solution:
    def purchasePlans(self, nums: List[int], target: int) -> int:
        nums.sort()
        res = 0
        left, right = 0, len(nums) - 1
        while left < right:  # 队列长度大于1
            if nums[left] + nums[right] <= target:
                res += right - left
                res %= MOD
                left += 1
            else:
                right -= 1
        return res


assert Solution().purchasePlans([2, 2, 1, 9], 10) == 4
