#            0 [0,5)
#           /       \
#          /         \
#        1 [0,3)      2 [3,4)
#       /    \
#      /      \
#    3 [0,1)   4[1,2)

from typing import List, Tuple


class DFSOrder:
    __slots__ = ("starts", "ends", "_n", "_tree", "_dfsId")

    def __init__(self, n: int, tree: List[List[int]], root=0) -> None:
        """dfs序

        Args:
            n (int): 树节点从0开始,根节点为0
            tree (Tree): 无向图邻接表

        1. 按照dfs序遍历k个结点形成的回路 每条边恰好经过两次
        """
        self.starts = [0] * n
        self.ends = [0] * n
        self._n = n
        self._tree = tree
        self._dfsId = 0
        self._dfs(root, -1)

    def querySub(self, root: int) -> Tuple[int, int]:
        """求子树映射到的区间

        Args:
            root (int): 根节点
        Returns:
            Tuple[int, int]: [start, end] 0 <= start < end <= n
        """
        return self.starts[root], self.ends[root] + 1

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


if __name__ == "__main__":
    N = 4
    edges = [[0, 1], [1, 2], [2, 3]]
    tree = [[] for _ in range(N)]
    for u, v in edges:
        tree[u].append(v)
        tree[v].append(u)
    D = DFSOrder(N, tree)
    print(D.querySub(1))
    print(D.querySub(2))
    print(D.querySub(3))
    print(D.queryId(1))
    print(D.queryId(2))
    print(D.queryId(3))
    print(D.querySub(3))
