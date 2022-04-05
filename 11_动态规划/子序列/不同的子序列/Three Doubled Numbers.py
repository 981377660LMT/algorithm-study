# return the number of triples i < j < k such that nums[i] * 2 = nums[j] and nums[j] * 2 = nums[k].
# 求数组中三元对的数量 (num,2*num,4*num)

# 可以用子序列dp


from collections import Counter
from typing import List


class Solution:
    def solve(self, nums: List[int]):
        dp1, dp2, dp3 = Counter(), Counter(), Counter()
        for num in nums:
            dp3[num] += dp2[num / 2]
            dp2[num] += dp1[num / 2]
            dp1[num] += 1
        return sum(dp3.values())


print(Solution().solve(nums=[1, 0, 2, 4, 4]))
