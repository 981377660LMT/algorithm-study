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
        self.depth = [0] * n  # 根节点深度为0
        self.parent = [0] * n  # 根节点父亲为-1

        self._n = n
        self._tree = tree

        self._dfs(root, -1, 0)
        self._bitlen = n.bit_length()
        self.dp = self._makeDp()  # !dp[i][j] 表示节点j的第2^i个祖先

    def queryLCA(self, root1: int, root2: int) -> int:
        """查询树节点两点的最近公共祖先"""
        if self.depth[root1] < self.depth[root2]:
            root1, root2 = root2, root1

        root1 = self._upToDepth(root1, self.depth[root2])
        if root1 == root2:
            return root1

        for i in range(self._bitlen - 1, -1, -1):
            if self.dp[i][root1] != self.dp[i][root2]:
                root1 = self.dp[i][root1]
                root2 = self.dp[i][root2]

        return self.dp[0][root1]

    def queryDist(self, root1: int, root2: int) -> int:
        """查询树节点两点间距离"""
        return self.depth[root1] + self.depth[root2] - 2 * self.depth[self.queryLCA(root1, root2)]

    def queryKthAncestor(self, root: int, k: int) -> int:
        """查询树节点root的第k个祖先,如果不存在这样的祖先节点,返回 -1"""
        bit = 0
        while k:
            if k & 1:
                root = self.dp[bit][root]
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

    def _dfs(self, start: int, startPre: int, startDep: int) -> None:
        """处理高度、父节点"""
        depth, parent, tree = self.depth, self.parent, self._tree
        stack = [(start, startPre, startDep)]
        while stack:
            cur, pre, dep = stack.pop()
            depth[cur], parent[cur] = dep, pre
            for next in tree[cur]:
                if next != pre:
                    stack.append((next, cur, dep + 1))

    def _makeDp(self) -> List[List[int]]:
        """nlogn预处理"""
        n, bitlen, parent = self._n, self._bitlen, self.parent
        dp = [[-1] * n for _ in range(bitlen)]
        for j in range(self._n):
            dp[0][j] = parent[j]
        for i in range(bitlen - 1):
            for j in range(n):
                if dp[i][j] == -1:
                    dp[i + 1][j] = -1
                else:
                    dp[i + 1][j] = dp[i][dp[i][j]]
        return dp

    def _upToDepth(self, root: int, toDepth: int) -> int:
        """从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点"""
        if toDepth >= self.depth[root]:
            return root
        for i in range(self._bitlen - 1, -1, -1):
            if (self.depth[root] - toDepth) & (1 << i):
                root = self.dp[i][root]
        return root


if __name__ == "__main__":
    n = 3
    adjList = [[1], [0, 2], [1]]
    lca = LCA(
        n,
        adjList,
        root=0,
    )
    # print(lca.queryLCA(1, 2))  # 1
    print(lca.queryLCA(0, 2))  # 2
