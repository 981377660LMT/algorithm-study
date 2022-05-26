class Node:
    def __init__(self, val):
        self.val = val
        self.left = None
        self.right = None
        self.parent = None


# 注意是二叉搜索树
# 给定一棵二叉搜索树和其中的一个节点 node ，找到该节点在树中的中序后继。如果节点没有中序后继，请返回 null 。
# 注意没有给出树
# 进阶：你能否在不访问任何结点的值的情况下解决问题?

# 1.有右子
# 则是右子树的最左下
# 2.无右子，需要往上找，且是第一个 左子--父 关系的父
class Solution:
    def inorderSuccessor(self, node: 'Node') -> 'Node':
        if node.right:
            right = node.right
            while right.left:
                right = right.left
            return right
        else:
            # node如果作为右节点，则需要向上继续找
            while node.parent and node.parent.right == node:
                node = node.parent
            return node.parent


print(Solution().inorderSuccessor())
