# 1 <= nums.length <= 105
# 统计和为正数的子数组有多少个，使用前缀和转换为计算左侧小于当前元素的个数的问题，可用有序集合或树状数组O(nlogn)解决
from typing import List
from sortedcontainers import SortedList

MOD = int(1e9 + 7)

# 此题可用树状数组/有序集合 计算左侧小于当前元素的个数
class Solution:
    def subarraysWithMoreZerosThanOnes(self, nums: List[int]) -> int:
        preSum = [0]
        for num in nums:
            preSum.append(preSum[-1] + (1 if num == 1 else -1))
        print(preSum)

        res = 0
        lis = SortedList()
        for s in preSum:
            index = lis.bisect_left(s)
            res += index
            lis.add(s)

        return res


print(Solution().subarraysWithMoreZerosThanOnes(nums=[0, 1, 1, 0, 1]))
