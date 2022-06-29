from collections import defaultdict
from typing import DefaultDict, Set, Tuple


class DFSOrder:
    def __init__(self, n: int, tree: DefaultDict[int, Set[int]]) -> None:
        """dfs序

        Args:
            n (int): 树节点从0开始,根节点为0
            tree (DefaultDict[int, Set[int]]): 无向图邻接表

        1. 按照dfs序遍历k个结点形成的回路 每条边恰好经过两次
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
        # !这里避免函数嵌套调用方法 self.queryRange 可以快1500ms
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


if __name__ == '__main__':
    N = 4
    edges = [[0, 1], [1, 2], [2, 3]]
    tree = defaultdict(set)
    for u, v in edges:
        tree[u].add(v)
        tree[v].add(u)
    dfsOrder = DFSOrder(N, tree)
    print(dfsOrder.queryRange(1))
    print(dfsOrder.queryRange(2))
    print(dfsOrder.queryRange(3))
    print(dfsOrder.queryId(1))
    print(dfsOrder.queryId(2))
    print(dfsOrder.queryId(3))
