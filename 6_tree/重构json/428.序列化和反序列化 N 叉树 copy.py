from typing import Optional


class Node(object):
    def __init__(self, val=None, children=None):
        self.val = val
        self.children = children


class Codec:
    # dfs 序列化和反序列化 N 叉树
    def serialize(self, root: 'Node') -> str:
        """节点值,子节点个数,节点值,子节点个数,..."""

        def dfs(root: Optional['Node']):
            if not root:
                return
            res.append(str(root.val))
            res.append(str(len(root.children)))
            for child in root.children:
                dfs(child)

        res = []
        dfs(root)
        return ",".join(res)

    def deserialize(self, data: str) -> Optional['Node']:
        def dfs() -> Node:
            root = Node(int(next(it)), [])
            childCount = int(next(it))
            for _ in range(childCount):
                root.children.append(dfs())
            return root

        if not data:
            return

        it = iter(data.split(","))
        return dfs()
