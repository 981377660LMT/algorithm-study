# 类似计算器，一个栈存操作符，一个栈存t/f

# 返回该式的运算结果。

# 遇到右括号，弹出1个操作符，再弹出操作数直至遇到左括号，
# 计算操作数与操作符的运算结果，将结果入栈
class Solution:
    def parseBoolExpr(self, expression: str) -> bool:
        nonOptStack, optStack = [], []

        for char in expression:
            if char == '|' or char == '&' or char == '!':
                optStack.append(char)
            elif char == 't':
                nonOptStack.append(True)
            elif char == 'f':
                nonOptStack.append(False)
            # 注意这里左括号进入nonOptStack，用于间隔t/f
            elif char == '(':
                nonOptStack.append('(')

            elif char == ')':

                # 与和或需要取出t/f,非直接append结果
                opt = optStack.pop()
                if opt == '&':

                    tOrF = nonOptStack.pop()
                    while nonOptStack[-1] != '(':
                        top = nonOptStack.pop()
                        tOrF = tOrF and top
                    # 弹出左括号
                    nonOptStack.pop()
                    nonOptStack.append(tOrF)

                elif opt == '|':
                    tOrF = nonOptStack.pop()
                    while nonOptStack[-1] != '(':
                        top = nonOptStack.pop()
                        tOrF = tOrF or top
                    # 弹出左括号
                    nonOptStack.pop()
                    nonOptStack.append(tOrF)

                elif opt == '!':
                    tOrF = nonOptStack.pop()
                    # 弹出左括号
                    nonOptStack.pop()
                    nonOptStack.append(not tOrF)

        return nonOptStack[0]


print(Solution().parseBoolExpr("&(t,f)"))
print(Solution().parseBoolExpr("|(f,t)"))
print(Solution().parseBoolExpr("|(&(t,f,t),!(t))"))
