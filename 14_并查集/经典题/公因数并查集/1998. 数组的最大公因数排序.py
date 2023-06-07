# https://leetcode.cn/problems/gcd-sort-of-an-array/
# 如果 gcd(nums[i], nums[j]) > 1 ，交换 nums[i] 和 nums[j] 的位置。
# 其中 gcd(nums[i], nums[j]) 是 nums[i] 和 nums[j] 的最大公因数。
# 如果能使用上述交换方式将 nums 按 非递减顺序 排列，返回 true ；否则，返回 false 。

from typing import List
from 埃氏筛和并查集 import EratosthenesSieve, UnionFindArray

E = EratosthenesSieve(int(1e5) + 10)


class Solution:
    def gcdSort(self, nums: List[int]) -> bool:
        max_ = max(nums, default=0)
        uf = UnionFindArray(max_ + 1)
        for num in nums:
            for p in E.getPrimeFactors(num):
                uf.union(num, p)

        target = sorted(nums)
        for a, b in zip(target, nums):
            if uf.find(a) != uf.find(b):
                return False
        return True
