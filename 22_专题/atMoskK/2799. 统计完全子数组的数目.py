# 完全子数组 ：
# !子数组中 不同 元素的数目等于整个数组不同元素的数目。
# 返回数组中 完全子数组 的数目。
# 对每个右端点求出左端点可以扩张的最远距离，然后累加即可。


from collections import defaultdict
from typing import List


class Solution:
    def countCompleteSubarrays(self, nums: List[int]) -> int:
        res, left, n = 0, 0, len(nums)
        target = len(set(nums))
        counter = defaultdict(int)
        for right in range(n):
            counter[nums[right]] += 1
            while left <= right and len(counter) == target:
                removed = nums[left]
                counter[removed] -= 1
                if counter[removed] == 0:
                    del counter[removed]
                left += 1
            res += left
        return res
