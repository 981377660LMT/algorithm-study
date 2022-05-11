# 类似计算器，一个栈存操作符，一个栈存t/f

# 返回该式的运算结果。

# 遇到右括号，弹出1个操作符，再弹出操作数直至遇到左括号，
# 计算操作数与操作符的运算结果，将结果入栈
class Solution:
    """这道题很特殊,左括号(哨兵)放在numStack里会容易处理一些"""

    def parseBoolExpr(self, expression: str) -> bool:
        numStack, optStack = [], []

        for char in expression:
            if char == '(':  # 注意这里左括号进入numStack，用作一段运算结束的哨兵
                numStack.append('(')
            elif char in ('|', '&', '!'):
                optStack.append(char)
            elif char in ('t', 'f'):
                numStack.append(True if char == 't' else False)
            elif char == ')':
                opt = optStack.pop()
                if opt == '&':
                    tOrF = numStack.pop()
                    while numStack[-1] != '(':
                        top = numStack.pop()
                        tOrF = tOrF and top
                    numStack.pop()  # 弹出左括号
                    numStack.append(tOrF)
                elif opt == '|':
                    tOrF = numStack.pop()
                    while numStack[-1] != '(':
                        top = numStack.pop()
                        tOrF = tOrF or top
                    numStack.pop()
                    numStack.append(tOrF)
                elif opt == '!':
                    tOrF = numStack.pop()
                    numStack.pop()
                    numStack.append(not tOrF)

        return numStack[0]


print(Solution().parseBoolExpr("&(t,f)"))
print(Solution().parseBoolExpr("|(f,t)"))
print(Solution().parseBoolExpr("|(&(t,f,t),!(t))"))
