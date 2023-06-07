# https://leetcode.cn/problems/graph-connectivity-with-threshold/
# x 和 y 的两座城市直接连通的前提是：
# !x 和 y 的公因数中，至少有一个 严格大于 某个阈值 threshold
# 给你两个整数 n 和 threshold ，以及一个待查询数组，
# 请你判断每个查询 queries[i] = [ai, bi] 指向的城市 ai 和 bi 是否连通
# 2 <= n <= 104
# 0 <= threshold <= n
# 注意threshold为 0 的情况(导致公因数不是质数)，需要单独处理

from typing import List
from 埃氏筛和并查集 import getFactors, UnionFindArray


class Solution:
    def areConnected(self, n: int, threshold: int, queries: List[List[int]]) -> List[bool]:
        """每个数与大于threshold的因子合并."""
        if threshold == 0:
            return [True] * len(queries)

        uf = UnionFindArray(n + 1)
        for num in range(threshold + 1, n + 1):
            for f in getFactors(num):
                if f > threshold:
                    uf.union(num, f)

        return [uf.isConnected(a, b) for a, b in queries]
