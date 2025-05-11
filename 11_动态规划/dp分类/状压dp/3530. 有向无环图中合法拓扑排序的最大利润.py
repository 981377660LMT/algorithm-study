# https://leetcode.cn/problems/maximum-profit-from-valid-topological-order-in-dag/solutions/3663054/pai-lie-xing-zhuang-ya-dpcong-ji-yi-hua-z67rp/
# !将每个节点的得分乘以其在拓扑排序中的位置，然后求和，得到的值称为 利润。
# 请返回在所有合法拓扑排序中可获得的 最大利润 。
# n<=22.

from typing import List
from functools import lru_cache


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxProfit(self, n: int, edges: List[List[int]], score: List[int]) -> int:
        if not edges:
            score.sort()
            return sum(i * s for i, s in enumerate(score, 1))

        pre = [0] * n
        for x, y in edges:
            pre[y] |= 1 << x

        @lru_cache(None)
        def dfs(mask: int) -> int:
            res = 0
            pos = mask.bit_count() + 1
            # 枚举未学习的课程，且该课程的前置课程都已学习
            for i, p in enumerate(pre):
                if mask >> i & 1 == 0 and (p & mask) == p:
                    res = max2(res, dfs(mask | (1 << i)) + pos * score[i])
            return res

        res = dfs(0)
        dfs.cache_clear()
        return res
