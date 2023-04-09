from bisect import bisect_right
from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。现有一个长度等于 nums.length 的数组 arr 。对于满足 nums[j] == nums[i] 且 j != i 的所有 j ，arr[i] 等于所有 |i - j| 之和。如果不存在这样的 j ，则令 arr[i] 等于 0 。


def calDistSum(nums: List[int], k: int, preSum: List[int]) -> int:
    """有序数组所有点到x=k的距离之和

    排序+二分+前缀和 O(logn)
    """
    pos = bisect_right(nums, k)
    leftSum = k * pos - preSum[pos]
    rightSum = preSum[-1] - preSum[pos] - k * (len(nums) - pos)
    return leftSum + rightSum


# 返回数组 arr 。
class Solution:
    def distance(self, nums: List[int]) -> List[int]:
        mp = defaultdict(list)
        for i in range(len(nums)):
            mp[nums[i]].append(i)
        res = [0] * len(nums)
        preSum = defaultdict(list)
        for k, v in mp.items():
            v.sort()
            preSum[k] = [0] + list(accumulate(v))
        for i in range(len(nums)):
            res[i] = calDistSum(mp[nums[i]], i, preSum[nums[i]])
        return res


# nums = [1,3,1,1,2]
print(Solution().distance([1, 3, 1, 1, 2]))
