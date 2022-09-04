# 找出环形数组中的连续子数组 和最接近某个数
# 一个环形数组，分成两半，求两半和的差的绝对值最小是多少？
# 环形数组分类讨论，分成两半,只有2种情况
# [AAAAAA BBBBBB]
# [BBB AAAAAA BBB]
# 枚举A
# !必然有其中一个区间是不跨越分界点的，所以求一个绝对值最接近一半的就行了，
# 如果都是正数的话双指针就行（刚好小于和刚好大于的都验一下），
# 如果有负数的话，前缀和放红黑树里面比较简单，复杂度O(nlogn)，复杂度O(n)

# !x+y=sum 求abs(x-y)最小值 即abs(2*子数组和-sum)最小值

from bisect import bisect_right
from typing import List
from sortedcontainers import SortedList

INF = int(1e20)


class Solution:
    def minDifference1(self, nums: List[int]) -> int:
        """都是正整数
        两半和的差的绝对值最小
        """
        target = sum(nums)
        nums = [num * 2 for num in nums]

        left, curSum = 0, 0
        res = INF
        for right in range(len(nums)):
            curSum += nums[right]
            res = min(res, abs(curSum - target))
            while left <= right and curSum > target:
                curSum -= nums[left]
                left += 1
                res = min(res, abs(curSum - target))

        return res

    def minDifference2(self, nums: List[int]) -> int:
        """有负数
        两半和的差的绝对值最小
        abs(2*子数组和-sum)

        1438.220. 存在重复元素 III
        """
        target = sum(nums)
        nums = [num * 2 for num in nums]

        preSum = SortedList()
        res, curSum = INF, 0
        for num in nums:
            curSum += num
            pos = bisect_right(preSum, curSum - target)
            if pos < len(preSum):
                # 注意这里是减去curSum - preSum[pos]
                res = min(res, abs((curSum - preSum[pos]) - target))
            if pos - 1 >= 0:
                res = min(res, abs((curSum - preSum[pos - 1]) - target))
            preSum.add(curSum)

        return res


print(Solution().minDifference1([1, 2, 3, 4]))
print(Solution().minDifference1([10, 2, 8, 3]))
print(Solution().minDifference2([1, 2, 3, 4]))
print(Solution().minDifference2([10, 2, 8, 3]))
