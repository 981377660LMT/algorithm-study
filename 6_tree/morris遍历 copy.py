class Tree:
    def __init__(self, val, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def solve(self, root):
        left = right = None
        last_node = Tree(-float("inf"))

        for node in self.morris_inorder(root):
            if node.val < last_node.val:
                left = left or last_node
                right = node
            last_node = node

        left.val, right.val = right.val, left.val

        return root

    @staticmethod
    def morris_inorder(node):
        temp = None
        while node:
            if node.left:
                temp = node.left
                while temp.right and temp.right != node:
                    temp = temp.right
                if temp.right:
                    temp.right = None
                    yield node
                    node = node.right
                else:
                    temp.right = node
                    node = node.left
            else:
                yield node
                node = node.right
