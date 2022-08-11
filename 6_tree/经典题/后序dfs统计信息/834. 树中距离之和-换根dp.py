from collections import defaultdict
from typing import Callable, List

# 834. 树中距离之和-换根dp
class Solution:
    def sumOfDistancesInTree(self, n: int, edges: List[List[int]]) -> List[int]:
        def dfs(cur: int, parent_: int, depth_: int) -> None:
            parent[cur] = parent_
            depth[cur] = depth_
            for next in adjMap[cur]:
                if next == parent_:
                    continue
                dfs(next, cur, depth_ + 1)
                subTreeCount[cur] += subTreeCount[next]

        def getRes(cur: int, parent: int) -> None:
            for next in adjMap[cur]:
                if next == parent:
                    continue
                # 注意这里都是 subTreeCount[next]
                res[next] = res[cur] - subTreeCount[next] + (n - subTreeCount[next])
                getRes(next, cur)

        depth = [-1] * n
        parent = [-1] * n
        subTreeCount = [1] * n

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        dfs(0, -1, 0)

        res = [0] * n
        res[0] = sum(depth)
        getRes(0, -1)
        return res


print(Solution().sumOfDistancesInTree(n=6, edges=[[0, 1], [0, 2], [2, 3], [2, 4], [2, 5]]))

#######################################
def rerooting(
    n: int,
    adjList: List[List[int]],
    init: Callable[[], int],
    child: Callable[[int], int],
    merge: Callable[[int, int], int],
):
    """换根dp框架 https://atcoder.jp/contests/dp/submissions/33369486

    以每个结点作为根 求某种性质

    Args:
        n: 顶点数
        adjList: 树 编号从1-n 0表示虚拟根节点
        init: 初始化根节点的值
        node: 子结点贡献值
        merge: 子节点与父节点的合并方式
    """

    stack = [(1, 0)]
    order = []
    while stack:
        pair = stack.pop()
        cur, pre = pair
        order.append(pair)
        for next in adjList[cur]:
            if next == pre:
                continue
            stack.append((next, cur))

    dp1 = [0] * (n + 1)
    ls = [0] * (n + 1)
    rs = [0] * (n + 1)

    for cur, pre in reversed(order):  # dfs序
        nexts = adjList[cur]
        v = init()
        for next in nexts:
            if next == pre:
                continue
            ls[next] = v
            v = merge(v, child(dp1[next]))
        v = init()
        for next in reversed(nexts):
            if next == pre:
                continue
            rs[next] = v
            v = merge(v, child(dp1[next]))
        dp1[cur] = v

    dp2 = [0] * (n + 1)
    res = [0] * (n + 1)
    res[1] = dp1[1]
    for i in range(1, n):
        cur, pre = order[i]
        dp2[cur] = v = merge(merge(ls[cur], rs[cur]), child(dp2[pre]))
        res[cur] = merge(child(v), dp1[cur])
    return res


class Solution2:
    def sumOfDistancesInTree(self, n: int, edges: List[List[int]]) -> List[int]:

        adjList = [[] for _ in range(n + 1)]
        for u, v in edges:
            u, v = u + 1, v + 1
            adjList[u].append(v)
            adjList[v].append(u)

        res = rerooting(n, adjList, lambda: 0, lambda x: x, lambda x, y: x + y + 1)
        return res[1:]


print(Solution2().sumOfDistancesInTree(n=6, edges=[[0, 1], [0, 2], [2, 3], [2, 4], [2, 5]]))
