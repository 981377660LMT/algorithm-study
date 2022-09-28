"""
给数组元素添加正负号,最小化`子数组和的绝对值的最大值`
1 <= nums.length <= 1000
1 <= nums[i] <= 1000


1. 参考 1749. 任意子数组和的绝对值的最大值,一个结论是：
任意子数组和的绝对值的最大值等于 max(preSum) - min(preSum)

2. 考虑二分答案mid,check 函数为给数组元素添加正负号后 max(preSum) - min(preSum) <=mid 是否成立

3. 记 dp[i][j] 为给数组元素添加正负号后,数组前 i 项的前缀和能取到 j

4. 将前缀和的最小值平移到 0,那么 dp 转移过程中前缀和不能超出 [0,mid] 这个区间；
从 [0,mid] 枚举所有可能的出发点,如果能完成 n 次转移,说明子数组和的绝对值的最大值不超过 mid
"""

from typing import List


class Solution:
    def unSuitability(self, nums: List[int]) -> int:
        """二分答案 + dp O(nUlogU)"""

        # def check(mid: int) -> bool:
        #     """
        #     给数组元素添加正负号后,max(preSum) - min(preSum) <= mid 是否成立
        #     """
        #     dp = set(range(mid + 1))
        #     for num in nums:
        #         ndp = set()
        #         for pre in dp:
        #             if pre + num <= mid:
        #                 ndp.add(pre + num)
        #             if pre - num >= 0:
        #                 ndp.add(pre - num)
        #         dp = ndp
        #     return len(dp) > 0

        def check(mid: int) -> bool:
            """
            给数组元素添加正负号后,max(preSum) - min(preSum) <= mid 是否成立
            """
            mask = (1 << (mid + 1)) - 1
            dp = mask
            for num in nums:
                dp = ((dp << num) | (dp >> num)) & mask
            return dp != 0

        left, right = 0, max(nums) * 2
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


print(Solution().unSuitability(nums=[5, 3, 7]))
print(Solution().unSuitability(nums=[20, 10]))
print(Solution().unSuitability(nums=[10, 20]))
