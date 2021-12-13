from typing import List
from collections import deque, defaultdict

# 给你数组 bombs ，请你返回在引爆 一个 炸弹的前提下，最多 能引爆的炸弹数目。

# 有向图
class Solution:
    def maximumDetonation(self, bombs: List[List[int]]) -> int:
        def bfs(cur: int):
            res = 0
            # ----bfs+记忆化
            queue = deque()
            visited = [False] * n
            queue.append(cur)
            visited[cur] = True
            while queue:
                x = queue.popleft()
                res += 1
                for y in adjMap[x]:
                    if not visited[y]:
                        visited[y] = True
                        queue.append(y)
            return res

        n = len(bombs)

        # --------建立有向图
        adjMap = defaultdict(list)
        for i in range(n):
            xi, yi, ri = bombs[i]
            for j in range(n):
                if i != j:
                    xj, yj, _ = bombs[j]
                    if (xi - xj) ** 2 + (yi - yj) ** 2 <= ri ** 2:
                        adjMap[i].append(j)

        # --------计算res
        ans = 0
        for i in range(n):
            ans = max(ans, bfs(i))
        return ans


print(Solution().maximumDetonation(bombs=[[2, 1, 3], [6, 1, 4]]))
