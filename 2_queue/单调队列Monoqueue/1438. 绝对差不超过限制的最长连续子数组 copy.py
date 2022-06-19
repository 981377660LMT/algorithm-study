from typing import List
from sortedcontainers import SortedList
from collections import deque

# 我们的目的是快速让一组数据有序，那就寻找一个内部是有序的数据结构
# 时间复杂度：O(N*log(N))
# 每个元素遍历一次，新元素插入红黑树的调整时间为 O(log(N))
# Python 的 SortedList 可以达到此目的。Java 可用 TreeMap，C++ 可用 multiset 代替。


# 1.自动有序的数据结构O(nlogn)
# 维护滑动窗口内的顺序性
class Solution:
    def longestSubarray(self, nums: List[int], limit: int) -> int:
        arr = SortedList()
        left, right = 0, 0
        res = 0

        for right in range(len(nums)):
            arr.add(nums[right])
            while arr[-1] - arr[0] > limit:
                arr.remove(nums[left])
                left += 1
            res = max(res, right - left + 1)

        return res

    def longestSubarray2(self, nums: List[int], limit: int) -> int:
        maxQueue, minQueue = deque(), deque()
        res, left, right = 0, 0, 0
        while right < len(nums):
            cur = nums[right]
            while minQueue and minQueue[-1] > cur:
                minQueue.pop()
            minQueue.append(cur)
            while maxQueue and maxQueue[-1] < cur:
                maxQueue.pop()
            maxQueue.append(cur)
            while maxQueue[0] - minQueue[0] > limit:
                maxQueue[0] == nums[left] and maxQueue.popleft()
                minQueue[0] == nums[left] and minQueue.popleft()
                left += 1
            res = max(res, right - left + 1)
            right += 1
        return res


print(Solution().longestSubarray([8, 2, 4, 7], 4))
print(Solution().longestSubarray2([8, 2, 4, 7], 4))

