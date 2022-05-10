from typing import Optional


class TreeNode:
    def __init__(self, x: int):
        self.val = x
        self.left: Optional['TreeNode'] = None
        self.right: Optional['TreeNode'] = None


class Codec:
    def serialize(self, root: TreeNode) -> str:
        def dfs(root: Optional['TreeNode']) -> None:
            if not root:
                return
            res.append(str(root.val))
            res.append('1' if root.left else '0')
            res.append('1' if root.right else '0')
            dfs(root.left)
            dfs(root.right)

        res = []
        dfs(root)
        return ",".join(res)

    def deserialize(self, data: str) -> Optional['TreeNode']:
        def dfs() -> TreeNode:
            root = TreeNode(int(next(it)))
            hasLeft, hasRight = next(it) == '1', next(it) == '1'
            if hasLeft:
                root.left = dfs()
            if hasRight:
                root.right = dfs()
            return root

        if not data:
            return
        it = iter(data.split(","))
        return dfs()
