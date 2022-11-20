from math import comb
from typing import List
from collections import Counter


# 给你一个下标从 0 开始的正整数数组 nums 。请你找出并统计满足下述条件的三元组 (i, j, k) 的数目：

# 0 <= i < j < k < nums.length
# nums[i]、nums[j] 和 nums[k] 两两不同 。
# 换句话说：nums[i] != nums[j]、nums[i] != nums[k] 且 nums[j] != nums[k] 。
# 返回满足上述条件三元组的数目。

# !6241. 数组中不等三元组的个数

# !1.nlogn/O(n)解法
# !三元组肯定是枚举中间那个数


class Solution:
    def unequalTriplets(self, nums: List[int]) -> int:
        """哈希表O(n)"""
        n, res = len(nums), 0
        counter = Counter(nums)
        smaller, bigger = 0, n
        for key in sorted(counter):  # !这里不sorted也可以(对称性)
            bigger -= counter[key]
            res += smaller * bigger * counter[key]
            smaller += counter[key]
        return res


print(Solution().unequalTriplets(nums=[4, 4, 2, 4, 3]))
