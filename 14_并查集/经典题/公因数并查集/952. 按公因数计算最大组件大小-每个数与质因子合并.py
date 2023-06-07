# https://leetcode.cn/problems/largest-component-size-by-common-factor/
# 只有当 A[i] 和 A[j] 共用一个大于 1 的公因数时，A[i] 和 A[j] 之间才有一条边。
# 只有当 nums[i] 和 nums[j] 共用`一个`大于 1 的公因数时，nums[i] 和 nums[j]之间才有一条边。
# 返回图中最大连通组件的大小。
# 1 <= A.length <= 2e4
# !1 <= A[i] <= 1e5

from collections import defaultdict
from typing import List
from 埃氏筛和并查集 import EratosthenesSieve, UnionFindArray

E = EratosthenesSieve(int(1e5) + 10)


class Solution:
    def largestComponentSize(self, nums: List[int]) -> int:
        """每个数与大于1的质因子合并(从上往下合并)"""
        max_ = max(nums, default=0)
        uf = UnionFindArray(max_ + 1)
        for num in nums:
            for factor in E.getPrimeFactors(num):
                uf.union(num, factor)

        group = defaultdict(int)
        for num in nums:
            group[uf.find(num)] += 1
        return max(group.values(), default=0)
