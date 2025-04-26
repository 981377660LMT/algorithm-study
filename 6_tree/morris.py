class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def morris_inorder(root):
    """
    Morris 中序遍历生成器，无需额外栈或递归，
    逐个 yield 出访问到的节点（按中序顺序）。
    完成后，能保证树结构被恢复。
    """
    curr = root
    while curr:
        if curr.left:
            # 找到左子树的最右节点
            pred = curr.left
            while pred.right and pred.right is not curr:
                pred = pred.right
            if pred.right is None:
                # 第一次访问：建立线索，进入左子树
                pred.right = curr
                curr = curr.left
            else:
                # 第二次访问：移除线索，访问 curr，然后进入右子树
                pred.right = None
                yield curr
                curr = curr.right
        else:
            # 无左子树，直接访问，再进入右子树
            yield curr
            curr = curr.right


if __name__ == "__main__":
    INF = int(1e18)

    class Solution:
        def recoverTree(self, root: TreeNode) -> None:
            """
            在常数额外空间下，恢复仅有两个节点被交换的 BST。
            """
            first = second = None
            prev = TreeNode(-INF)

            for node in morris_inorder(root):
                if prev.val > node.val:
                    if first is None:
                        first = prev
                    second = node
                prev = node

            first.val, second.val = second.val, first.val  # type: ignore
