# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def rob(self, root: 'TreeNode') -> int:
        def dfs(root: 'TreeNode'):
            if not root:
                return [0, 0]
            res = [0, root.val]  # 不选 选
            for next in (root.left, root.right):
                noNext, hasNext = dfs(next)
                res[0] += max(noNext, hasNext)
                res[1] += noNext
            return res

        return max(dfs(root))
