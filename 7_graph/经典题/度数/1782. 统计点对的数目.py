from typing import List
from collections import Counter

# 2 <= n <= 2 * 104


# 第 j 个查询的答案是满足如下条件的点对 (a, b) 的数目：

# a < b
# cnt 是与 a 或者 b 相连的边的数目，且 cnt 严格大于 queries[j] 。

# 请注意，图中可能会有 重复边 。(需要用(min,max)统计数量去重)


# https://leetcode.com/problems/count-pairs-of-nodes/discuss/1096740/C%2B%2BJavaPython3-Two-Problems-O(q-*-(n-%2B-e))


class Solution:
    def countPairs(self, n: int, edges: List[List[int]], queries: List[int]) -> List[int]:
        degree, res = [0] * (n + 1), [0] * len(queries)
        # 统计重复边
        overlap = Counter((min(u, v), max(u, v)) for u, v in edges)
        for u, v in edges:
            degree[u] += 1
            degree[v] += 1

        sortedDegree = sorted(degree)

        # 排序的两数之和
        for index, query in enumerate(queries):
            left, right = 1, n
            while left < right:
                if query < sortedDegree[left] + sortedDegree[right]:
                    res[index] += right - left
                    right -= 1
                else:
                    left += 1

            for (u, v), count in overlap.items():
                # 重复边：度数重复算了count 次
                if query < degree[u] + degree[v] and query >= degree[u] + degree[v] - count:
                    res[index] -= 1

        return res


print(Solution().countPairs(n=4, edges=[[1, 2], [2, 4], [1, 3], [2, 3], [2, 1]], queries=[2, 3]))
# 输出：[6,5]
# 解释：每个点对中，与至少一个点相连的边的数目如上图所示。


# todo 太难了

