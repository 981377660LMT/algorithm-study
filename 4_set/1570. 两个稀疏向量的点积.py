# https://leetcode.cn/problems/dot-product-of-two-sparse-vectors/solutions/2381771/liang-ge-xi-shu-xiang-liang-de-dian-ji-b-2ljd/

from typing import List


class SparseVector:
    __slots__ = ("data",)

    def __init__(self, nums: List[int]):
        self.data = {i: v for i, v in enumerate(nums) if v}

    def dotProduct(self, vec: "SparseVector") -> int:
        # 遍历更稀疏的一侧
        if len(self.data) > len(vec.data):
            return vec.dotProduct(self)
        return sum(val * vec.data.get(idx, 0) for idx, val in self.data.items())


class SparseVector2:
    __slots__ = ("pairs",)

    def __init__(self, nums: List[int]):
        self.pairs = []
        for i, v in enumerate(nums):
            if v != 0:
                self.pairs.append([i, v])

    def dotProduct(self, vec: "SparseVector2") -> int:
        res = 0
        p, q = 0, 0
        pairs1, pairs2 = self.pairs, vec.pairs
        while p < len(pairs1) and q < len(pairs2):
            if pairs1[p][0] == pairs2[q][0]:
                res += pairs1[p][1] * pairs2[q][1]
                p += 1
                q += 1
            elif pairs1[p][0] < pairs2[q][0]:
                p += 1
            else:
                q += 1

        return res
