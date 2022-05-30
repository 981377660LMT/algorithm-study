from typing import List, Optional


class TreeNode:
    def __init__(self, val, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

    def __repr__(self):
        return 'TreeNode({})'.format(self.val)


def deserializeNode(data: List[Optional[int]]) -> Optional['TreeNode']:
    """反序列化二叉树"""

    def genNode(val: Optional[int]) -> Optional['TreeNode']:
        if val is None:
            return None
        return TreeNode(val)

    if not data:
        return None
    root = genNode(data.pop(0))
    queue = [root]

    while queue:
        cur = queue.pop(0)
        if cur:
            cur.left = genNode(data.pop(0) if data else None)
            cur.right = genNode(data.pop(0) if data else None)
            if cur.left:
                queue.append(cur.left)
            if cur.right:
                queue.append(cur.right)
    return root


def drawTree(root: Optional['TreeNode']):
    """画出二叉树"""
    import turtle

    def getHeight(root: Optional['TreeNode']) -> int:
        return 1 + max(getHeight(root.left), getHeight(root.right)) if root else -1

    def jumpto(x: float, y: float) -> None:
        t.penup()
        t.goto(x, y)
        t.pendown()

    def draw(node: Optional['TreeNode'], x: float, y: float, dx: float):
        if node:
            t.goto(x, y)
            jumpto(x, y - 20)
            t.write(node.val, align='center', font=('Arial', 12, 'normal'))
            draw(node.left, x - dx, y - 60, dx / 2)
            jumpto(x, y - 20)
            draw(node.right, x + dx, y - 60, dx / 2)

    t = turtle.Turtle()
    t.speed(0)
    turtle.delay(0)
    h = getHeight(root)
    jumpto(0, 30 * h)
    draw(root, 0, 30 * h, 40 * h)
    t.hideturtle()
    turtle.mainloop()


if __name__ == '__main__':
    print(deserializeNode([1, 2, 3, None, None, 4, 5]))
    drawTree(
        deserializeNode(
            [2, 1, 3, 0, 7, 9, 1, 2, None, 1, 0, None, None, 8, 8, None, None, None, None, 7]
        )
    )

