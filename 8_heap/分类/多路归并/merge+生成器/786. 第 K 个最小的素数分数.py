from heapq import merge
from itertools import islice
from typing import List

# 时间复杂度: O(klogn)
# n<=1000
# 其中k最大1e6


class Solution:
    def kthSmallestPrimeFraction(self, nums: List[int], k: int) -> List[int]:
        """求第k小的arr[i] / arr[j]"""
        n, nums = len(nums), sorted(nums)
        gen = lambda i: ((nums[i] / nums[j], nums[i], nums[j]) for j in range(n - 1, i, -1))  # 递增
        allGen = [gen(i) for i in range(n)]
        iterable = merge(*allGen, key=lambda x: x[0])
        res = next(islice(iterable, k - 1, None))
        return [res[1], res[2]]


print(Solution().kthSmallestPrimeFraction(nums=[1, 2, 3, 5], k=3))
