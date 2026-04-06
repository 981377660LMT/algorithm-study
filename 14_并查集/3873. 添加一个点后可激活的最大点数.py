# 3873. 添加一个点后可激活的最大点数
# https://leetcode.cn/problems/maximum-points-activated-with-one-addition/description/
# 找到两个包含点数之和最大的连通块，通过一个新点将它们连起来

from collections import defaultdict
from heapq import nlargest
from UnionFind import UnionFindMap


OFFSET = int(1e10)


class Solution:
    def maxActivated(self, points: list[list[int]]) -> int:
        uf = UnionFindMap()
        for r, c in points:
            uf.union(r, c + OFFSET)
        groupSize = defaultdict(int)
        for r, c in points:
            groupSize[uf.find(r)] += 1
        return sum(nlargest(2, groupSize.values())) + 1
