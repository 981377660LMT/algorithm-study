from typing import Optional, Set


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 其中有且只有一个无效节点，它的右子节点错误地指向了与其在同一层且在其右侧的一个其他节点。
# 给定一棵这样的问题二叉树的根节点 root ，将该无效节点及其所有子节点移除（除被错误指向的节点外），然后返回新二叉树的根结点。

# 因为错误节点的右指针指向右侧的某一节点，所以利用先访问右孩子的先序遍历(根右左)，
# 当访问到错误节点时，其右指针指向的孩子肯定已经被访问过了，
# 根据这一点判断当前节点是否为错误节点
class Solution:
    def correctBinaryTree(self, root: TreeNode) -> TreeNode:
        def dfs(root: TreeNode, visited: Set[TreeNode]) -> TreeNode:
            if not root:
                return None

            # 待删除的结点
            if root.right in visited:
                return None

            visited.add(root)
            root.right = dfs(root.right, visited)
            root.left = dfs(root.left, visited)
            return root

        visited = set()
        return dfs(root, visited)

