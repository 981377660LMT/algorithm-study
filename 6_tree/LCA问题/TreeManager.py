from collections import defaultdict
from math import floor, log2
from typing import DefaultDict, Iterable, List, Mapping, Sequence, Tuple, Union


ListTree = Sequence[Iterable[int]]  # Sequence = Iterable + __getitem__
DictTree = Mapping[int, Iterable[int]]
Tree = Union[ListTree, DictTree]


class TreeManager:
    def __init__(self, n: int, tree: Tree, /, *, root: int, useLCA: bool) -> None:
        """查询 DFS 序 / LCA / 距离 / 路径

        - LCA :
        `nlogn` 预处理
        `logn`查询两点的LCA

        Args:
            n (int): 树节点编号 默认 0 ~ n-1
            tree (Tree): 树
            root (int): 根节点
            useLCA (bool): 是否使用倍增求LCA
        """
        self.depth = defaultdict(lambda: -1)  # 深度 一般用深度来排序结点 根节点深度为0
        self.parent = defaultdict(lambda: -1)  # 父节点 一般用来上跳查找路径 根节点父亲为-1
        self.start = [0] * n
        self.end = [0] * n

        self._n = n
        self._dfsId = 1
        self._useLCA = useLCA
        self._tree = tree

        self._dfs(root, -1, 0)
        if useLCA:
            self._BITLEN = floor(log2(n)) + 1
            self._fa = self._makeDp(self.parent)

    def queryLCA(self, root1: int, root2: int) -> int:
        """`logn` 查询"""
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

    def queryRange(self, root: int) -> Tuple[int, int]:
        """求子树映射到的区间

        Args:
            root (int): 根节点
        Returns:
            Tuple[int, int]: [start, end] 1 <= start <= end <= n
        """
        return self.start[root], self.end[root]

    def queryId(self, root: int) -> int:
        """求root自身的dfsId

        Args:
            root (int): 根节点
        Returns:
            int: id  1 <= id <= n
        """
        return self.end[root]

    def isAncestor(self, *, root: int, child: int) -> bool:
        """判断root是否是child的祖先

        Args:
            root (int): 根节点
            child (int): 子节点

        应用:枚举边时给树的边定向
        ```
        if not D.isAncestor(e[0], e[1]):
            e[0], e[1] = e[1], e[0]
        ```
        """
        left1, right1 = self.start[root], self.end[root]
        left2, right2 = self.start[child], self.end[child]
        return left1 <= left2 <= right2 <= right1

    def isLink(self, nodes: List[int]) -> bool:
        """判断结点是否组成从根节点出发的链

        TODO 可能会有问题
        https://zhuanlan.zhihu.com/p/540022071
        """
        nodes = sorted(nodes, key=lambda x: self.depth[x])
        for i in range(len(nodes) - 1):
            if not self.isAncestor(root=nodes[i], child=nodes[i + 1]):
                return False
        return True

    def isSimplePath(self, nodes: List[int]) -> bool:
        """判断结点是否组成一条简单路径(起点+一个拐点+终点)

        TODO 可能会有问题
        https://zhuanlan.zhihu.com/p/540022071
        """
        if len(nodes) <= 2:
            return True

        nodes = sorted(nodes, key=lambda x: self.depth[x])
        start = nodes[-1]
        anotherBranch = []
        for node in nodes[:-1]:
            if not self.isAncestor(root=node, child=start):
                anotherBranch.append(node)

        if not anotherBranch:  # !一条链
            return True

        anotherBranch.sort(key=lambda x: self.depth[x])
        end = anotherBranch[-1]
        uTurn = self.queryLCA(start, end)  # 拐点

        for node in nodes:
            if not self.isAncestor(root=uTurn, child=node):  # 拐点不是结点的祖先
                return False
            if not (
                self.isAncestor(root=node, child=start)  # 结点是起点的祖先
                or self.isAncestor(root=node, child=end)  # 结点是终点的祖先
            ):
                return False
        return True

    def _dfs(self, cur: int, pre: int, dep: int) -> None:
        """处理高度、父节点、dfs序信息"""
        self.depth[cur], self.parent[cur] = dep, pre
        self.start[cur] = self._dfsId
        for next in self._tree[cur]:
            if next == pre:
                continue
            self._dfs(next, cur, dep + 1)
        self.end[cur] = self._dfsId
        self._dfsId += 1

    def _makeDp(self, parent: DefaultDict[int, int]) -> List[List[int]]:
        """nlogn预处理"""
        dp = [[0] * self._BITLEN for _ in range(self._n)]
        for i in range(self._n):
            dp[i][0] = parent[i]
        for j in range(self._BITLEN - 1):
            for i in range(self._n):
                if dp[i][j] == -1:
                    dp[i][j + 1] = -1
                else:
                    dp[i][j + 1] = dp[dp[i][j]][j]
        return dp


if __name__ == "__main__":
    adjList1 = [[0, 1], [0, 2], [0, 3], [1, 4], [1, 5]]
    adjList2 = [set([0, 1]), set([0, 2]), set([0, 3]), set([1, 4]), set([1, 5])]
    adjMap = defaultdict(set)
    TreeManager(
        6,
        adjMap,
        useLCA=True,
        root=0,
    )
