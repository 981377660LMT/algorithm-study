from typing import List
from sortedcontainers import SortedList

# 我们的目的是快速让一组数据有序，那就寻找一个内部是有序的数据结构
# 时间复杂度：O(N*log(N))
# 每个元素遍历一次，新元素插入红黑树的调整时间为 O(log(N))
# Python 的 SortedList 可以达到此目的。Java 可用 TreeMap，C++ 可用 multiset 代替。


# 1.自动有序的数据结构O(nlogn)
# 维护滑动窗口内的顺序性
class Solution:
    def longestSubarray(self, nums: List[int], limit: int) -> int:
        arr = SortedList()
        l = 0
        r = 0
        res = 0

        while r < len(nums):
            arr.add(nums[r])
            while arr[-1] - arr[0] > limit:
                arr.remove(nums[l])
                l += 1
            res = max(res, r - l + 1)
            r += 1

        return res


print(Solution().longestSubarray([8, 2, 4, 7], 4))

