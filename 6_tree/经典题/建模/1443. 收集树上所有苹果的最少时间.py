# 通过树上的一条边，需要花费 1 秒钟
# 你从 节点 0 出发，请你返回最少需要多少秒，可以收集到所有苹果，并回到节点 0 。
# 1 <= n <= 10^5
# 无向边edges

# 总结：看走了多少条边
# 1.将需要收集的苹果的祖先节点的hasApple状态自下而上标记为True
# 2.遍历无向树edges中各个边，如果构成边的两点的hasApple状态都为True，说明需要经过这条边，res+=1，而收集的过程中每条边需要走两次，所以2*res即为所求
# https://leetcode-cn.com/problems/minimum-time-to-collect-all-apples-in-a-tree/solution/python3-zi-di-xiang-shang-dfs-by-yim-6-aub7/
from collections import defaultdict
from typing import List


class Solution:
    def minTime(self, n: int, edges: List[List[int]], hasApple: List[bool]) -> int:
        def dfs(cur: int, pre: int) -> None:
            """后序遍历 把苹果路径传上来"""
            for next in adjMap[cur]:
                if next == pre:
                    continue
                dfs(next, cur)
                if hasApple[next]:
                    hasApple[cur] = True

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        dfs(0, -1)
        return 2 * sum((hasApple[u] and hasApple[v]) for u, v in edges)

