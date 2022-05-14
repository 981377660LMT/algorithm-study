# 求最长的前缀，使得移除一个元素之后，剩下所有元素freq相等

# during iteration
# 1. decrement the biggest count or
# 2. decrement the smallest count
# 移除最多或者最少的


from typing import List
from collections import Counter
from sortedcontainers import SortedList


class Solution:
    def maxEqualFreq(self, nums: List[int]) -> int:
        """有序集合维护有序的频率"""
        n = len(nums)
        counter = Counter(nums)
        sl = SortedList(counter.values())

        for i in range(n - 1, -1, -1):
            # 一个
            if len(sl) <= 1:
                return i + 1
            # 删最少的 1 3 3 3
            if sl[0] == 1 and sl[1] == sl[-1]:
                return i + 1
            # 删最多的 3 3 3 4
            if sl[0] == sl[-2] and sl[-1] == sl[-2] + 1:
                return i + 1

            num = nums[i]
            sl.discard(counter[num])
            if counter[num] - 1 > 0:
                sl.add(counter[num] - 1)
            counter[num] -= 1
            if counter[num] == 0:
                del counter[num]

        return 0


print(Solution().maxEqualFreq([1, 1, 1, 2, 2, 3]))
print(Solution().maxEqualFreq(nums=[2, 2, 1, 1, 5, 3, 3, 5]))
print(Solution().maxEqualFreq(nums=[1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5]))
