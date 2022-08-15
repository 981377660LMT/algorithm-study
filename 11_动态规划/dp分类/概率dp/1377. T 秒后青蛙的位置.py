from typing import List

# 一棵由 n 个顶点组成的无向树，顶点编号从 1 到 n
# 青蛙从 顶点 1 开始起跳。规则如下：
# 如果青蛙可以跳到多个不同顶点，那么它跳到其中任意一个顶点上的机率都相同。
# 如果青蛙不能跳到任何未访问过的顶点上，那么它每次跳跃都会停留在原地。
# 返回青蛙在 t 秒后位于目标顶点 target 上的概率。

# dfs后序遍历 获取到下面传上来的概率
# !此节点处的概率就是 当前贡献/分支数
# 无路可走或者到点时 判断cur是否等于target
class Solution:
    def frogPosition(self, n: int, edges: List[List[int]], t: int, target: int) -> float:
        if n == 1:
            return 1

        adjlist = [[] for _ in range(n + 1)]
        for u, v in edges:
            adjlist[u].append(v)
            adjlist[v].append(u)

        visited = [False] * (n + 1)

        def dfs(cur: int, time: int) -> float:
            if time == t or (len(adjlist[cur]) - (cur != 1)) == 0:
                return cur == target
            visited[cur] = True
            res = sum(dfs(next, time + 1) for next in adjlist[cur] if not visited[next])
            return res / (len(adjlist[cur]) - (cur != 1))  # - (cur != 1) 表示减去来的路，起点没有来的路所以不减

        return dfs(1, 0)


print(Solution().frogPosition(7, [[1, 2], [1, 3], [1, 7], [2, 4], [2, 6], [3, 5]], 2, 4))
