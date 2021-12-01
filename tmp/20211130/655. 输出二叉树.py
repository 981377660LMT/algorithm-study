# 在一个 m*n 的二维字符串数组中输出二叉树
# 行数 m 应当等于给定二叉树的高度。
# 列数 n 应当总是奇数。


from typing import List


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def printTree(self, root: TreeNode) -> List[List[str]]:
        if not root:
            return [[""]]

        def depth(root):
            if not root:
                return 0
            return max(depth(root.left), depth(root.right)) + 1

        def gen(root: TreeNode, depth: int, pos: int) -> None:
            self.res[-depth - 1][pos] = str(root.val)
            if root.left:
                gen(root.left, depth - 1, pos - 2 ** (depth - 1))
            if root.right:
                gen(root.right, depth - 1, pos + 2 ** (depth - 1))

        d = depth(root)
        self.res = [[""] * (2 ** d - 1) for _ in range(d)]
        gen(root, d - 1, 2 ** (d - 1) - 1)
        return self.res


# 输入:
#      1
#     /
#    2
# 输出:
# [["", "1", ""],
#  ["2", "", ""]]

# 输入:
#      1
#     / \
#    2   3
#     \
#      4
# 输出:
# [["", "", "", "1", "", "", ""],
#  ["", "2", "", "", "", "3", ""],
#  ["", "", "4", "", "", "", ""]]

