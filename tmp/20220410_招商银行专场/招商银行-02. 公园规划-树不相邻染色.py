from collections import defaultdict
from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)

# 树中边的两个端点颜色不能相同
# 并且与同一个点相连的所有点，颜色必须互不相同
# 请返回这些花坛最少需要几种花。
# 招商银行-02. 公园规划-树不相邻染色


class Solution:
    def numFlowers(self, roads: List[List[int]]) -> int:
        degree = defaultdict(int)
        for u, v in roads:
            degree[v] += 1
            degree[u] += 1
        return max(degree.values()) + 1


print(Solution().numFlowers(roads=[[0, 1], [1, 3], [1, 2]]))
print(Solution().numFlowers(roads=[[0, 1], [0, 2], [1, 3], [2, 5], [3, 6], [5, 4]]))
