# 1 <= expression.length <= 1e5
# expression only contains '1','0','&','|','(', and ')'
class Solution:
    def minOperationsToFlip(self, expression: str) -> int:
        """操作符栈+操作数栈

        https://leetcode-cn.com/problems/minimum-cost-to-change-the-final-value-of-expression/solution/zhan-dong-tai-gui-hua-by-lucifer1004-7bsn/
        个人感觉不太好理解
        """
        optStack = []
        numStack = []
        for char in expression:
            if char in '01)':
                if char == '0':
                    numStack.append((0, 1))  # 变为0的操作数，变为1的操作数
                elif char == '1':
                    numStack.append((1, 0))
                else:
                    assert optStack[-1] == '('
                    optStack.pop()

                # 操作
                if optStack and optStack[-1] != '(':
                    opt = optStack.pop()
                    left0, left1 = numStack.pop()
                    right0, right1 = numStack.pop()
                    if opt == '&':
                        # 只要1边为0就变0，两边都为1或者边操作符为|
                        numStack.append(
                            (min(right0, left0), min(right1 + left1, 1 + min(right1, left1)))
                        )
                    elif opt == '|':
                        numStack.append(
                            (min(right0 + left0, 1 + min(right0, left0)), min(right1, left1))
                        )
            else:
                optStack.append(char)

        # 所有操作执行完毕后，我们的操作数栈将只包含一个元素。这个元素必定包含一个零值（对应于表达式原本的值）和一个非零值。而这个非零值就是我们要寻找的答案。
        return max(numStack[0])


print(Solution().minOperationsToFlip(expression="1&(0|1)"))

# 输出：1
# 解释：我们可以将 "1&(0|1)" 变成 "1&(0&1)" ，执行的操作为将一个 '|' 变成一个 '&' ，执行了 1 次操作。
# 新表达式的值为 0 。

