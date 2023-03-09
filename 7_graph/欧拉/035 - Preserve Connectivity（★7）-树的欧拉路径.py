"""!连接k个树结点需要保留的最小边数

n<=1e5

多个点 按照dfs序(树的欧拉路径) 求两两间的距离

树的欧拉序
"""

from collections import defaultdict
import sys
from math import floor, log2
from typing import DefaultDict, List, Set
from collections import defaultdict
from typing import DefaultDict, Set, Tuple


class _DFSOrder:
    def __init__(self, n: int, tree: DefaultDict[int, Set[int]]) -> None:
        """dfs序

        Args:
            n (int): 树节点从0开始,根节点为0
            tree (DefaultDict[int, Set[int]]): 无向图邻接表
        """
        self.starts = [0] * n
        self.ends = [0] * n
        self._n = n
        self._tree = tree
        self._dfsId = 1
        self._dfs(0, -1)

    def queryRange(self, root: int) -> Tuple[int, int]:
        """求子树映射到的区间

        Args:
            root (int): 根节点
        Returns:
            Tuple[int, int]: [start, end] 1 <= start <= end <= n
        """
        return self.starts[root], self.ends[root]

    def queryId(self, root: int) -> int:
        """求root自身的dfsId

        Args:
            root (int): 根节点
        Returns:
            int: id  1 <= id <= n
        """
        return self.ends[root]

    def isAncestor(self, root: int, child: int) -> bool:
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
        left1, right1 = self.starts[root], self.ends[root]
        left2, right2 = self.starts[child], self.ends[child]
        return left1 <= left2 <= right2 <= right1

    def _dfs(self, cur: int, pre: int) -> None:
        self.starts[cur] = self._dfsId
        for next in self._tree[cur]:
            if next == pre:
                continue
            self._dfs(next, cur)
        self.ends[cur] = self._dfsId
        self._dfsId += 1


#######################################################################
from collections import defaultdict, deque
from math import floor, log2
from typing import Iterable, List, Mapping, Sequence, Union


ListTree = Sequence[Iterable[int]]  # Sequence = Iterable + __getitem__
DictTree = Mapping[int, Iterable[int]]
Tree = Union[ListTree, DictTree]


class _LCA:
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

        self._bfs(root, -1, 0)
        self._bitlen = floor(log2(n)) + 1
        self._fa = self._makeDp(self.parent)

    def queryLCA(self, root1: int, root2: int) -> int:
        """查询树节点两点的最近公共祖先"""
        if self.depth[root1] < self.depth[root2]:
            root1, root2 = root2, root1

        for i in range(self._bitlen - 1, -1, -1):
            if self.depth[self._fa[root1][i]] >= self.depth[root2]:
                root1 = self._fa[root1][i]

        if root1 == root2:
            return root1

        for i in range(self._bitlen - 1, -1, -1):
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

    def _bfs(self, start: int, startPre: int, startDep: int) -> None:
        """处理高度、父节点"""
        queue = deque([(start, startPre, startDep)])
        while queue:
            cur, pre, dep = queue.popleft()
            self.depth[cur], self.parent[cur] = dep, pre
            for next in self._tree[cur]:
                if next != pre:
                    queue.append((next, cur, dep + 1))

    def _makeDp(self, parent: List[int]) -> List[List[int]]:
        """nlogn预处理"""
        fa = [[-1] * self._bitlen for _ in range(self._n)]
        for i in range(self._n):
            fa[i][0] = parent[i]
        for j in range(self._bitlen - 1):
            for i in range(self._n):
                if fa[i][j] == -1:
                    fa[i][j + 1] = -1
                else:
                    fa[i][j + 1] = fa[fa[i][j]][j]
        return fa


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())
adjMap = defaultdict(set)
for _ in range(n - 1):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjMap[u].add(v)
    adjMap[v].add(u)

D = _DFSOrder(n, adjMap)
L = _LCA(n, adjMap)

q = int(input())
for _ in range(q):
    k, *nodes = list(map(int, input().split()))
    nodes = [int(v) - 1 for v in nodes]
    nodes.sort(key=lambda x: D.queryId(x))
    res = 0
    for pre, cur in zip(nodes, nodes[1:] + [nodes[0]]):
        res += L.queryDist(pre, cur)
    print(res // 2)
