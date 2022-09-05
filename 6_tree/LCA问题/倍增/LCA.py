from collections import defaultdict
from math import floor, log2
from typing import Iterable, List, Mapping, Sequence, Union


ListTree = Sequence[Iterable[int]]  # Sequence = Iterable + __getitem__
DictTree = Mapping[int, Iterable[int]]
Tree = Union[ListTree, DictTree]


class LCA:
    def __init__(self, n: int, tree: Tree, root: int) -> None:
        """倍增查询LCA

        `nlogn` 预处理
        `logn`查询

        Args:
            n (int): 树节点编号 默认 0 ~ n-1
            tree (Tree): 树
            root (int): 根节点
        """
        self.depth = [-1] * (n + 1)  # 用深度来排序结点 根节点深度为0 最后一个结点-1表示虚拟结点
        self.parent = [-1] * (n + 1)  # 父节点用来上跳查找路径 根节点父亲为-1

        self._n = n
        self._tree = tree

        self._dfs(root, -1, 0)
        self._BITLEN = floor(log2(n)) + 1
        self._fa = self._makeDp(self.parent)

    def queryLCA(self, root1: int, root2: int) -> int:
        """查询树节点两点的最近公共祖先"""
        if self.depth[root1] < self.depth[root2]:
            root1, root2 = root2, root1

        for i in range(self._BITLEN - 1, -1, -1):
            if self.depth[self._fa[root1][i]] >= self.depth[root2]:
                root1 = self._fa[root1][i]

        if root1 == root2:
            return root1

        for i in range(self._BITLEN - 1, -1, -1):
            if self._fa[root1][i] != self._fa[root2][i]:
                root1 = self._fa[root1][i]
                root2 = self._fa[root2][i]

        return self._fa[root1][0]

    def queryDist(self, root1: int, root2: int) -> int:
        """查询树节点两点间距离"""
        return self.depth[root1] + self.depth[root2] - 2 * self.depth[self.queryLCA(root1, root2)]

    def queryKthAncestor(self, root: int, k: int) -> int:
        """查询树节点root的第k个祖先,如果不存在这样的祖先节点,返回 -1"""
        bit = 0
        while k:
            if k & 1:
                root = self._fa[root][bit]
                if root == -1:
                    return -1
            bit += 1
            k //= 2
        return root

    # def isLink(self, nodes: List[int]) -> bool:
    #     """判断结点是否组成从根节点出发的链

    #     https://zhuanlan.zhihu.com/p/540022071
    #     """
    #     nodes = sorted(nodes, key=lambda x: self.depth[x])
    #     for i in range(len(nodes) - 1):
    #         if not self.isAncestor(root=nodes[i], child=nodes[i + 1]):
    #             return False
    #     return True

    # def isSimplePath(self, nodes: List[int]) -> bool:
    #     """判断结点是否组成一条简单路径(起点+一个拐点+终点)

    #     https://zhuanlan.zhihu.com/p/540022071
    #     """
    #     if len(nodes) <= 2:
    #         return True

    #     nodes = sorted(nodes, key=lambda x: self.depth[x])
    #     start = nodes[-1]
    #     anotherBranch = []
    #     for node in nodes[:-1]:
    #         if not self.isAncestor(root=node, child=start):
    #             anotherBranch.append(node)

    #     if not anotherBranch:  # !一条链
    #         return True

    #     anotherBranch.sort(key=lambda x: self.depth[x])
    #     end = anotherBranch[-1]
    #     uTurn = self.queryLCA(start, end)  # 拐点

    #     for node in nodes:
    #         if not self.isAncestor(root=uTurn, child=node):  # 拐点不是结点的祖先
    #             return False
    #         if not (
    #             self.isAncestor(root=node, child=start)  # 结点是起点的祖先
    #             or self.isAncestor(root=node, child=end)  # 结点是终点的祖先
    #         ):
    #             return False
    #     return True

    def _dfs(self, cur: int, pre: int, dep: int) -> None:
        """处理高度、父节点"""
        self.depth[cur], self.parent[cur] = dep, pre
        for next in self._tree[cur]:
            if next == pre:
                continue
            self._dfs(next, cur, dep + 1)

    def _makeDp(self, parent: List[int]) -> List[List[int]]:
        """nlogn预处理"""
        fa = [[-1] * self._BITLEN for _ in range(self._n)]
        for i in range(self._n):
            fa[i][0] = parent[i]
        for j in range(self._BITLEN - 1):
            for i in range(self._n):
                if fa[i][j] == -1:
                    fa[i][j + 1] = -1
                else:
                    fa[i][j + 1] = fa[fa[i][j]][j]
        return fa


if __name__ == "__main__":
    adjList1 = [[0, 1], [0, 2], [0, 3], [1, 4], [1, 5]]
    adjList2 = [set([0, 1]), set([0, 2]), set([0, 3]), set([1, 4]), set([1, 5])]
    adjMap = defaultdict(set)
    LCA(
        6,
        adjMap,
        root=0,
    )
