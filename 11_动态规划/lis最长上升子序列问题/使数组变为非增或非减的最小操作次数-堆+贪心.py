# Make Array Non-decreasing or Non-increasing
# 每次操作可以使每个数加1或者减1


# 1 <= nums.length <= 1e5
# 0 <= nums[i] <= 1e9

# 显然dp[i][num]

from heapq import heappop, heappush, heappushpop, heapreplace
from typing import List

# https://codeforces.com/blog/entry/47821  slope trick
class Solution:
    def convertArray(self, nums: List[int]) -> int:
        def helper(nums: List[int]) -> int:
            """变为不减数组的最小操作次数
            
            如果num比前面的数小，那么就把前面的最大数变小
            """
            res, pq = 0, []  # 大根堆
            for num in nums:
                preMax = -pq[0]
                if preMax > num:  # 之前的最大值是这个值的瓶颈
                    res += preMax - num
                    heappushpop(pq, -num)
                heappush(pq, -num)
            return res

        return min(helper(nums), helper(nums[::-1]))


print(Solution().convertArray(nums=[3, 2, 4, 5, 0]))
print(Solution().convertArray(nums=[3, 1, 2, 1]))
print(Solution().convertArray([11, 11, 13, 8, 18, 19, 20, 7, 16, 3]))


# def heappushpop(heap, item):
#     """Fast version of a heappush followed by a heappop."""
#     if heap and heap[0] < item:
#         item, heap[0] = heap[0], item
#         _siftup(heap, 0)
#     return item


# def heapreplace(heap, item):
#     """Pop and return the current smallest value, and add the new item.

#     This is more efficient than heappop() followed by heappush(), and can be
#     more appropriate when using a fixed-size heap.  Note that the value
#     returned may be larger than item!  That constrains reasonable uses of
#     this routine unless written as part of a conditional replacement:

#         if item > heap[0]:
#             item = heapreplace(heap, item)
#     """
#     returnitem = heap[0]  # raises appropriate IndexError if heap is empty
#     heap[0] = item
#     _siftup(heap, 0)
#     return returnitem
