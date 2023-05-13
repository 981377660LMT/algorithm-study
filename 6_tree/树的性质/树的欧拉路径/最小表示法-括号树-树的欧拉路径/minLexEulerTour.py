# 字典序最小欧拉路径/括号树的最小表示法
# O(nlogn) 双端队列优化字符串的拼接.

# !0表示进入子树,1表示返回父节点.


from collections import deque
from typing import Deque, List


def minLexEulerTour01(tree: List[List[int]]) -> List[int]:
    """返回树的长为2*n的01欧拉序列,其中n为节点数.
    0表示进入子树,1表示返回父节点.
    """

    def dfs(cur: int, pre: int) -> Deque[int]:
        sub = sorted([dfs(next, cur) for next in tree[cur] if next != pre])
        res = deque()
        for d in sub:
            if len(res) > len(d):
                while d:
                    res.append(d.popleft())
            else:
                res, d = d, res
                while d:
                    res.appendleft(d.pop())
        res.appendleft(0)
        res.append(1)
        return res

    res = dfs(0, -1)
    return list(res)


def minLexEulerTourDep(tree: List[List[int]]) -> List[int]:
    """返回树的长为2*n-1的表示深度的欧拉序列,其中n为节点数."""

    def dfs(cur: int, pre: int, dep: int) -> Deque[int]:
        sub = sorted([dfs(next, cur, dep + 1) for next in tree[cur] if next != pre])
        res = deque([dep])  # !进入
        for d in sub:
            if len(res) > len(d):
                while d:
                    res.append(d.popleft())
            else:
                res, d = d, res
                while d:
                    res.appendleft(d.pop())
            res.append(dep)  # !回溯
        return res

    res = dfs(0, -1, 0)
    return list(res)


if __name__ == "__main__":
    n = 4
    edges = [[0, 1], [1, 2], [0, 3]]
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    print(minLexEulerTour01(adjList))
    print(minLexEulerTourDep(adjList))
