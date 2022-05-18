from typing import List
from sortedcontainers import SortedList

# 1 <= nums.length <= 105
# -104 <= nums[i] <= 104
class Solution:
    def countSmaller(self, nums: List[int]):
        n = len(nums)
        res = [0] * n
        visited = SortedList()
        # 用i或者deque会比append加反转快1000ms
        for i in range(n - 1, -1, -1):
            # 每次遍历开始比较i与j的关系,二分法看看是第几位
            smaller = visited.bisect_left(nums[i])
            res[i] = smaller
            visited.add(nums[i])

        return res

