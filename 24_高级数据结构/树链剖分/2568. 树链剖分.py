# 详见算法竞赛进阶指南\GoDS (Go Data Structures)\src\graph\heavylightdecomposition\heavylightdecomposition.go
# !deprecated

# 给定一棵树，树中包含 n 个节点（编号 1∼n），其中第 i 个节点的权值为 ai。
# 初始时，1 号节点为树的根节点。

from collections import defaultdict


# 1. 预处理所有节点的 重儿子 以及子树内 节点的数量 和 每个节点的 父节点
def dfs1(cur: int, pre: int, dep: int) -> None:
    depths[cur] = dep
    parents[cur] = pre
    subsizes[cur] = 1
    for next in adjMap[cur]:
        if next == pre:
            continue
        dfs1(next, cur, dep + 1)
        subsizes[cur] += subsizes[next]
        # 重儿子是子树节点最多的儿子
        if subsizes[next] > subsizes[heavySons[cur]]:
            heavySons[cur] = next


# 2. 树链剖分，找出每个节点所属 重链 的 顶点，dfs序的编号，并记录每个序号对应的树结点值
def dfs2(cur: int, heavyStart: int) -> None:
    """dfs2做剖分

    heavyStart: 重链顶点的编号
    """
    global id
    dfsIds[cur] = id
    valueById[id] = values[cur]
    heavyTops[cur] = heavyStart
    id += 1
    if heavySons[cur] == 0:  # 叶节点结束
        return
    dfs2(heavySons[cur], heavyStart)  # 先看重儿子重链剖分
    # 再处理轻儿子
    for next in adjMap[cur]:
        if (next == heavySons[cur]) or (next == parents[cur]):
            continue
        dfs2(next, next)  # 轻儿子的重链顶点就是他自己


def updateRange(root1: int, root2: int, delta: int) -> None:
    while heavyTops[root1] != heavyTops[root2]:  # 向上爬找到相同重链
        if depths[heavyTops[root1]] < depths[heavyTops[root2]]:
            root1, root2 = root2, root1  # root1的重链顶点要更深
            # dfs序原因，上面节点的id必然小于下面节点的id
        bit.add(dfsIds[heavyTops[root1]], dfsIds[root1], delta)
        root1 = parents[heavyTops[root1]]
    if depths[root1] < depths[root2]:
        root1, root2 = root2, root1
    bit.add(dfsIds[root2], dfsIds[root1], delta)  # 在同一重链中，处理剩余区间


def updateRoot(root: int, delta: int) -> None:
    """子树全部加上delta 由于dfs序的原因,可以利用子树节点个数直接找到区间"""
    bit.add(dfsIds[root], dfsIds[root] + subsizes[root] - 1, delta)


def queryRange(root1: int, root2: int) -> int:
    res = 0
    while heavyTops[root1] != heavyTops[root2]:  # 向上爬找到相同重链
        if depths[heavyTops[root1]] < depths[heavyTops[root2]]:
            root1, root2 = root2, root1  # root1要更深
            # dfs序原因，上面节点的id必然小于下面节点的id
        res += bit.query(dfsIds[heavyTops[root1]], dfsIds[root1])
        root1 = parents[heavyTops[root1]]
    if depths[root1] < depths[root2]:
        root1, root2 = root2, root1
    res += bit.query(dfsIds[root2], dfsIds[root1])  # 在同一重链中，处理剩余区间
    return res


def queryRoot(root: int) -> int:
    return bit.query(dfsIds[root], dfsIds[root] + subsizes[root] - 1)


class BIT:
    def __init__(self, n: int):
        self.size = n
        self._tree1 = defaultdict(int)
        self._tree2 = defaultdict(int)

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left, right]加delta"""
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left, right]的和"""
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")

        rawIndex = index
        while index <= self.size:
            self._tree1[index] += delta
            self._tree2[index] += (rawIndex - 1) * delta
            index += index & -index

    def _query(self, index: int) -> int:
        if index > self.size:
            index = self.size

        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1[index] - self._tree2[index]
            index -= index & -index
        return res


n = int(input())
values = [0] + list(map(int, input().split()))

adjMap = defaultdict(set)
for _ in range(n - 1):
    u, v = map(int, input().split())
    adjMap[u].add(v)
    adjMap[v].add(u)

depths = [0] * (n + 1)
parents = [0] * (n + 1)
subsizes = [0] * (n + 1)
heavySons = [0] * (n + 1)  # 每个点的重儿子
heavyTops = [0] * (n + 1)  # 每个点的重链顶点
dfsIds = [0] * (n + 1)
valueById = [0] * (n + 1)
id = 1

dfs1(1, -1, 1)
dfs2(1, 1)
bit = BIT(n + 10)
for i in range(1, n + 1):
    bit.add(i, i, valueById[i])

m = int(input())
for _ in range(m):
    opt, *args = map(int, input().split())
    if opt == 1:
        # 修改路径权值
        u, v, delta = args
        #
        updateRange(u, v, delta)
    elif opt == 2:
        root, delta = args
        # 将以节点 u 为根的子树上的所有节点的权值增加 k。
        updateRoot(root, delta)
    elif opt == 3:
        u, v = args
        # 询问节点 u 和节点 v 之间路径上的所有节点（包括这两个节点）的权值和。
        print(queryRange(u, v))
    elif opt == 4:
        root = args[0]
        # 询问以节点 u 为根的子树上的所有节点的权值和。
        print(queryRoot(root))
# 5
# 1 3 7 4 5
# 1 3
# 1 4
# 1 5
# 2 3
# 5
# 1 3 4 3
# 3 5 4


# !如果支持断开边 那就需要LCT(link-cut tree)了
