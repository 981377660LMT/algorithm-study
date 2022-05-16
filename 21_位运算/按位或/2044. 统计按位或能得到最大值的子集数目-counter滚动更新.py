from typing import Counter, List


class Solution:
    def countMaxOrSubsets(self, nums: List[int]) -> int:
        """统计按位或能得到最大值的子集数目"""

        # counter滚动更新
        dp = Counter({0: 1})
        for num in nums:
            for key, count in list(dp.items()):
                dp[key | num] += count
        return dp[max(dp)]

    def countMaxOrSubsets2(self, nums: List[int]) -> int:
        """统计按位或能得到最大值的子集数目
        
        递推求所有子集的按位或
        """

        dp = [0]
        for num in nums:
            ndp = []
            for pre in dp:
                ndp.append(pre | num)
            dp += ndp
        counter = Counter(dp)
        return counter[max(counter)]


print(Solution().countMaxOrSubsets(nums=[3, 2, 1, 5]))
# 输出：6
# 解释：子集按位或可能的最大值是 7 。有 6 个子集按位或可以得到 7 ：
# - [3,5]
# - [3,1,5]
# - [3,2,5]
# - [3,2,1,5]
# - [2,5]
# - [2,1,5]
