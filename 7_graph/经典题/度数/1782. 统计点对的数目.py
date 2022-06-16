from typing import List
from collections import Counter

# 2 <= n <= 2 * 1e4
# 1 <= queries.length <= 20

# 第 j 个查询的答案是满足如下条件的点对 (a, b) 的数目：

# a < b
# cnt 是与 a 或者 b 相连的边的数目，且 cnt `严格`大于 queries[j] 。

# !请注意，图中可能会有 重复边 。(需要用(min,max)统计数量去重)


# https://leetcode.com/problems/count-pairs-of-nodes/discuss/1096740/C%2B%2BJavaPython3-Two-Problems-O(q-*-(n-%2B-e))


class Solution:
    def countPairs(self, n: int, edges: List[List[int]], queries: List[int]) -> List[int]:
        """排序+双指针"""
        deg, res = [0] * (n + 1), [0] * len(queries)
        # 统计重复边
        edgeCounter = Counter((min(u, v), max(u, v)) for u, v in edges)
        for u, v in edges:
            deg[u] += 1
            deg[v] += 1

        # 排序的两数之和
        sortedDeg = sorted(deg)
        for qi, qv in enumerate(queries):
            p1, p2 = 1, n
            while p1 < p2:
                if qv < sortedDeg[p1] + sortedDeg[p2]:
                    res[qi] += p2 - p1
                    p2 -= 1
                else:
                    p1 += 1

            for (u, v), count in edgeCounter.items():
                # !重复边导致的假阳性：度数重复算了count 次
                if qv < deg[u] + deg[v] and qv >= deg[u] + deg[v] - count:
                    res[qi] -= 1

        return res


print(Solution().countPairs(n=4, edges=[[1, 2], [2, 4], [1, 3], [2, 3], [2, 1]], queries=[2, 3]))
# 输出：[6,5]
# 解释：每个点对中，与至少一个点相连的边的数目如上图所示。


# todo 太难了

