# !RMQ  子树区间前后缀最大值
# !看到子树 直接弄到区间上

# !给定q个查询,每个查询给定一个节点x,要求删除以x为根的子树,并求出删除后的二叉树的高度。
# 查询之间是独立的，所以在每个查询执行后，树会回到其 初始 状态。

from itertools import accumulate
from typing import List, Optional


INF = int(1e18)


class Solution:
    def treeQueries(self, root: Optional["TreeNode"], queries: List[int]) -> List[int]:
        def dfs(root: Optional["TreeNode"], dep: int) -> None:
            if not root:
                return
            nonlocal dfsId
            ins[root.val] = dfsId
            dfs(root.left, dep + 1)
            dfs(root.right, dep + 1)
            outs[root.val] = dfsId
            depth[dfsId] = dep
            dfsId += 1

        n = int(1e5)
        ins, outs, dfsId = [0] * (n + 10), [0] * (n + 10), 1
        depth = [0] * (n + 10)
        dfs(root, 0)

        preMax = list(accumulate(depth, max))
        sufMax = list(accumulate(depth[::-1], max))[::-1]
        res = [0] * len(queries)
        for i in range(len(queries)):
            queryRoot = queries[i]
            left, right = ins[queryRoot], outs[queryRoot]
            max1, max2 = preMax[left - 1], sufMax[right + 1]
            res[i] = max(max1, max2)  # type: ignore

        return res


if __name__ == "__main__":

    class TreeNode:
        def __init__(self, val, left=None, right=None):
            self.val = val
            self.left = left
            self.right = right

        def __repr__(self):
            return "TreeNode({})".format(self.val)

    def deserializeNode(data: List[Optional[int]]) -> Optional["TreeNode"]:
        """反序列化二叉树"""

        def genNode(val: Optional[int]) -> Optional["TreeNode"]:
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

    def drawTree(root: Optional["TreeNode"]):
        """画出二叉树"""
        import turtle

        def getHeight(root: Optional["TreeNode"]) -> int:
            return 1 + max(getHeight(root.left), getHeight(root.right)) if root else -1

        def jumpto(x: float, y: float) -> None:
            t.penup()
            t.goto(x, y)
            t.pendown()

        def draw(node: Optional["TreeNode"], x: float, y: float, dx: float):
            if node:
                t.goto(x, y)
                jumpto(x, y - 20)
                t.write(node.val, align="center", font=("Arial", 12, "normal"))
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

    node = deserializeNode([1, None, 5, 3, None, 2, 4])
    print(Solution().treeQueries(node, [3, 5, 4, 2, 4]))  # type: ignore
    drawTree(node)
# [1,0,3,3,3]
