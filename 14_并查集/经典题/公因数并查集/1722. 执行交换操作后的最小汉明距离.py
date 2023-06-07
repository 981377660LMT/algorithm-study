# https://leetcode.cn/problems/minimize-hamming-distance-after-swap-operations/
# 每个 allowedSwaps[i] = [ai, bi] 表示你可以交换数组 source 中下标为 ai 和 bi.
# 相同长度的两个数组 source 和 target 间的 汉明距离 是元素不同的下标数量.
# !在对数组 source 执行 任意 数量的交换操作后，返回 source 和 target 间的 最小汉明距离 。


# 1. 并查集获取帮派邻接表.
# 2. 计算每个连通块对应的source元素与target的交集.


from collections import Counter
from typing import List
from 埃氏筛和并查集 import UnionFindArray


class Solution:
    def minimumHammingDistance(
        self, source: List[int], target: List[int], allowedSwaps: List[List[int]]
    ) -> int:
        n = len(source)
        uf = UnionFindArray(n)
        for a, b in allowedSwaps:
            uf.union(a, b)

        samePair = 0
        for g in uf.getGroups().values():
            c1 = Counter(source[i] for i in g)
            c2 = Counter(target[i] for i in g)
            samePair += sum((c1 & c2).values())
        return n - samePair


# source = [1,2,3,4], target = [2,1,4,5], allowedSwaps = [[0,1],[2,3]]
print(Solution().minimumHammingDistance([1, 2, 3, 4], [2, 1, 4, 5], [[0, 1], [2, 3]]))
