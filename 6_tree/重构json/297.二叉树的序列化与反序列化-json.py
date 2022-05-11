import json
from typing import Optional


class TreeNode:
    def __init__(self, x: int):
        self.val = x
        self.left: Optional['TreeNode'] = None
        self.right: Optional['TreeNode'] = None


class Codec:
    def serialize(self, root: TreeNode) -> str:
        def dfs(root: Optional['TreeNode']):
            if not root:
                return
            return {'val': root.val, 'left': dfs(root.left), 'right': dfs(root.right)}

        return json.dumps(dfs(root))

    def deserialize(self, data: str) -> Optional['TreeNode']:
        def dfs(obj):
            if obj is None:
                return
            return TreeNode(obj['val'], dfs(obj['left']), dfs(obj['right']))

        return dfs(json.loads(data))
