# n 块石头放置在二维平面中的一些整数坐标点上。每个坐标点上最多只能有一块石头。
# !如果一块石头的 同行或者同列 上有其他石头存在，那么就可以移除这块石头。
# 给你一个长度为 n 的数组 stones ，
# 其中 stones[i] = [xi, yi] 表示第 i 块石头的位置，
# !返回 可以移除的石子 的最大数量。


# 可以将石子全部建立并查集的联系，并计算联通分量的个数。
# 答案就是 n - 联通分量的个数 (每个联通分量至少要保留一块石头)。
from itertools import combinations
from typing import List
from UnionFind import UnionFindArray


class Solution:
    def removeStones(self, stones: List[List[int]]) -> int:
        n = len(stones)
        uf = UnionFindArray(n)
        for i, j in combinations(range(n), 2):
            if stones[i][0] == stones[j][0] or stones[i][1] == stones[j][1]:
                uf.union(i, j)
        return n - uf.part


assert Solution().removeStones([[0, 0], [0, 1], [1, 0], [1, 2], [2, 1], [2, 2]]) == 5
