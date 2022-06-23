from typing import Optional, Tuple


class MyNode:
    def __init__(self, val='0', left: Optional['MyNode'] = None, right: Optional['MyNode'] = None):
        self.val = val
        self.left = left
        self.right = right


# 1 <= expression.length <= 1e5
# expression only contains '1','0','&','|','(', and ')'
# 可以将0变为1，1变为0，&变为|，|变为&，
# 求反转表达式的最少操作次数
class Solution:
    def minOperationsToFlip(self, expression: str) -> int:
        """根据中缀表达式构造二叉表达式树然后树上dp
        叶子结点表示操作数,非叶子结点表示运算符
        可以反转任意结点的值
        """

        def dfs(root: Optional['MyNode']) -> Tuple[int, int]:
            """返回(变为0的最小操作次数,变为1的最小操作次数)"""
            if not root:
                return int(1e20), int(1e20)
            if root.val.isdigit():
                return int(root.val == '1'), int(root.val == '0')

            left0, left1 = dfs(root.left)
            right0, right1 = dfs(root.right)
            res0, res1 = int(1e20), int(1e20)

            assert root.val in ('&', '|')
            if root.val == '&':
                res0 = min(res0, left0 + right0, left0 + right1, left1 + right0)
                res1 = min(res1, left1 + right1, left0 + right1 + 1, left1 + right0 + 1)
            else:
                res0 = min(res0, left0 + right0, left0 + right1 + 1, left1 + right0 + 1)
                res1 = min(res1, left1 + right1, left0 + right1, left1 + right0)

            return res0, res1

        weight = dict(zip('(&|', (0, 1, 1)))
        numStack, optStack = [], []
        expression += ')'
        for char in expression:
            if char == '(':
                optStack.append(char)
            elif char.isdigit():
                numStack.append(MyNode(char))
            elif char in '&|':
                while optStack and weight[optStack[-1]] >= weight[char]:
                    node2, node1 = numStack.pop(), numStack.pop()
                    numStack.append(MyNode(optStack.pop(), node1, node2))
                optStack.append(char)
            elif char == ')':
                while optStack and optStack[-1] != '(':
                    ndoe2, node1 = numStack.pop(), numStack.pop()
                    numStack.append(MyNode(optStack.pop(), node1, ndoe2))
                if optStack:
                    optStack.pop()

        root = numStack[0]
        # print(json.dumps(root, indent=2, default=lambda o: o.__dict__))
        # drawtree(root)
        # print(dfs(root))
        return max(dfs(root))


if __name__ == '__main__':

    def drawtree(root):
        def height(root):
            return 1 + max(height(root.left), height(root.right)) if root else -1

        def jumpto(x, y):
            t.penup()
            t.goto(x, y)
            t.pendown()

        def draw(node, x, y, dx):
            if node:
                t.goto(x, y)
                jumpto(x, y - 20)
                t.write(node.val, align='center', font=('Arial', 12, 'normal'))
                draw(node.left, x - dx, y - 60, dx / 2)
                jumpto(x, y - 20)
                draw(node.right, x + dx, y - 60, dx / 2)

        import turtle

        t = turtle.Turtle()
        t.speed(0)
        turtle.delay(0)
        h = height(root)
        jumpto(0, 30 * h)
        draw(root, 0, 30 * h, 40 * h)
        t.hideturtle()
        turtle.mainloop()

    # print(Solution().minOperationsToFlip(expression="1&(0|1)"))
    # 输出：1
    # 解释：我们可以将 "1&(0|1)" 变成 "1&(0&1)" ，执行的操作为将一个 '|' 变成一个 '&' ，执行了 1 次操作。
    # 新表达式的值为 0 。
    # print(Solution().minOperationsToFlip(expression="(0&0)&(0&0&0)"))
    print(Solution().minOperationsToFlip(expression="1|1|(0&0)&1"))
    # print(eval("1|1|(0&0)&1"))

