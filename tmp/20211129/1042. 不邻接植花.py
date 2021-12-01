from typing import List

# 四色定理
# 在每个花园中，你打算种下`四种花之一`,所有花园 最多 有 3 条路径可以进入或离开。
# 你需要为每个花园选择一种花，使得通过路径相连的任何两个花园中的花的种类互不相同。
# 以数组形式返回 任一 可行的方案作为答案 answer，其中 answer[i] 为在第 (i+1) 个花园中种植的花的种类

# 贪心即可，因为所有花园 最多 有 3 条路径可以进入或离开
class Solution:
    def gardenNoAdj(self, n: int, paths: List[List[int]]) -> List[int]:
        res = [0] * n
        adjList = [[] for _ in range(n)]
        for u, v in paths:
            adjList[u - 1].append(v - 1)
            adjList[v - 1].append(u - 1)
        for i in range(n):
            # 随机pop出一种
            res[i] = ({1, 2, 3, 4} - {res[j] for j in adjList[i]}).pop()

        return res


print(Solution().gardenNoAdj(n=3, paths=[[1, 2], [2, 3], [3, 1]]))
