from typing import List
from LCA import LCA


class TreeAncestor:
    def __init__(self, n: int, parent: List[int]):
        adjList = [[] for _ in range(n)]
        for cur, pre in enumerate(parent):
            if pre == -1:
                continue
            adjList[pre].append(cur)
            adjList[cur].append(pre)
        self._lca = LCA(n, adjList, root=0)

    def getKthAncestor(self, node: int, k: int) -> int:
        return self._lca.queryKthAncestor(node, k)


if __name__ == "__main__":
    treeAncestor = TreeAncestor(7, [-1, 0, 0, 1, 1, 2, 2])
    print(treeAncestor.getKthAncestor(3, 1))  # 返回 1 ，它是 3 的父节点
    print(treeAncestor.getKthAncestor(5, 2))  # 返回 0 ，它是 5 的祖父节点
    print(treeAncestor.getKthAncestor(6, 3))  # 返回 -1 ，因为不存在满足要求的祖先节点
