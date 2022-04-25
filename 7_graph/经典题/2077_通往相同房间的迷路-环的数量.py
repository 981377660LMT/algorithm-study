from typing import List
from collections import defaultdict
from itertools import combinations


# 求无向图中长度为 3 的不同环的数量
class Solution:
    def numberOfPaths(self, n: int, corridors: List[List[int]]) -> int:
        adjMap = defaultdict(set)
        for u, v in corridors:
            adjMap[u].add(v)
            adjMap[v].add(u)

        res = 0
        for p1, p2 in corridors:
            res += len(adjMap[p1] & adjMap[p2])
        return res

        # res = 0
        # # 对邻居两两判断
        # for p1 in range(1, n + 1):
        #     for p2, p3 in combinations(adj[p1], 2):
        #         if p2 in adj[p3]:
        #             res += 1
        # return res // 3


# Output: 2
print(Solution().numberOfPaths(n=5, corridors=[[1, 2], [5, 2], [4, 1], [2, 4], [3, 1], [3, 4]]))

# 把无向图建成有向图
# 由于走廊是双向的，我们直接遍历会重复计算，在每一个存在环的三元组中，
# 最小的节点左右路径肯定是最大和第二大的节点，所以我们只需遍历比当前大的节点，
# 这样问题就变成了有向图，就可以避免重复计算了，也减少了数据量。

