"""排序+并查集"""

from typing import List
from RangeUnionFind import UnionFindRange


class Solution:
    def validSubarraySize(self, nums: List[int], threshold: int) -> int:
        """并查集维护区间标记 从大到小遍历 把看过的区间串起来"""
        n = len(nums)
        uf = UnionFindRange(n + 10)
        Q = sorted(((num, i) for i, num in enumerate(nums)), reverse=True)  # 数组中的元素越大越好
        for num, i in Q:
            uf.union(i, i + 1)  # 向右连接
            length = uf.getSize(uf.find(i + 1)) - 1  # 串联区间的长度
            if num * length > threshold:
                return length
        return -1
