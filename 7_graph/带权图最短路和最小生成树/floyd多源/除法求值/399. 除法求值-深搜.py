from collections import defaultdict
from functools import lru_cache
from typing import Generator, List, Set

# 1 <= equations.length <= 20
# 0.0 < values[i] <= 20.0
# !a/b 为2 则a−>b连一条边，边权为2.0 ; b−>a连一条边，边权为1.0/2.0。
# 如果存在某个无法确定的答案，则用 -1.0 替代这个答案。
# 如果问题中出现了给定的已知条件中没有出现的字符串，也需要用 -1.0 替代这个答案
# !你可以假设除法运算中不会出现除数为 0 的情况，且不存在任何矛盾的结果


class Solution:
    def calcEquation(
        self, equations: List[List[str]], values: List[float], queries: List[List[str]]
    ) -> List[float]:
        def dfs(
            cur: str, curValue: float, target: str, visited: Set[str]
        ) -> Generator[float, None, None]:
            """求cur到target的(最短)路径长度 深搜宽搜都可以
            
            生成器写法简洁
            """
            if cur not in adjMap or target not in adjMap:
                yield -1.0
            if cur == target:
                yield curValue
            for next in adjMap[cur]:
                if next in visited:
                    continue
                visited.add(next)
                yield from dfs(next, curValue * adjMap[cur][next], target, visited)

        adjMap = defaultdict(lambda: defaultdict(lambda: 0.0))
        for (u, v), w in zip(equations, values):
            adjMap[u][v] = w
            adjMap[v][u] = 1 / w
        return [next(dfs(u, 1.0, v, set([u])), -1.0) for u, v in queries]


print(
    Solution().calcEquation(
        equations=[["a", "b"], ["b", "c"], ["bc", "cd"]],
        values=[1.5, 2.5, 5.0],
        queries=[["a", "c"], ["c", "b"], ["bc", "cd"], ["cd", "bc"]],
    )
)
