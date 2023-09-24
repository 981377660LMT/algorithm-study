from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 下标从 0 开始的数组 nums ，数组中的元素为 互不相同 的正整数。请你返回让 nums 成为递增数组的 最少右移 次数，如果无法得到递增数组，返回 -1 。


# 一次 右移 指的是同时对所有下标进行操作，将下标为 i 的元素移动到下标 (i + 1) % n 处。


class Solution:
    def minimumRightShifts(self, nums: List[int]) -> int:
        def check(nums, count):
            arr = nums[:]
            for i in range(len(nums)):
                arr[i] = nums[(i - count) % len(nums)]
            return sorted(arr) == arr

        res = -1
        for count in range(len(nums) + 2):
            if check(nums, count):
                res = count
                break
        return res
