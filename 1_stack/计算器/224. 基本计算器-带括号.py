import re
from math import nan
from operator import add, sub, mul, truediv, floordiv

WEIGHT = dict(zip('(+-*/', [nan, 1, 1, 2, 2]))
CAL = dict(zip('+-*/', [add, sub, mul, floordiv]))


class Solution:
    def calculate2(self, s: str) -> int:
        """不能用eval python的eval有层数限制200 而js没有
        
        SyntaxError: too many nested parentheses
        直接eval(s)的话会报MemoryError，考虑从内到外依次计算括号
        """
        while True:
            left = s.rfind('(')
            right = s.find(')', left)
            if left == -1 or right == -1:
                break
            s = s[:left] + str(eval(s[left : right + 1])) + s[right + 1 :]
        return eval(s)

    def calculate(self, s: str) -> int:
        """基本计算器-带括号"""
        numStack, optStack = [], []

        if s[0] == '-':  # (+digit (-digit 以及开头的-digit
            s = '0' + s
        s = s.replace('(-', '(0-').replace('(+', '(0+')
        tokens = [v.strip() for v in re.split(r'([\(\+\-\*\/\)])', s) if v.strip()]
        tokens.append(')')

        for char in tokens:
            if char.isdigit():
                numStack.append(int(char))
            else:
                if char != ')':
                    while optStack and WEIGHT[optStack[-1]] >= WEIGHT[char]:
                        num2, num1 = numStack.pop(), numStack.pop()
                        numStack.append(CAL[optStack.pop()](num1, num2))
                    optStack.append(char)
                else:
                    while optStack and optStack[-1] != '(':
                        num2, num1 = numStack.pop(), numStack.pop()
                        numStack.append(CAL[optStack.pop()](num1, num2))
                    if optStack:
                        optStack.pop()

        return numStack[0]


if __name__ == '__main__':
    print(Solution().calculate("1 + 1"))
    print(Solution().calculate("   (  3 ) "))
    assert Solution().calculate("(1+(4+5+2)-3)+(6+8)") == 23

    print(re.split(r'(\+|\-|\*|\/)', '1+1-22'))
    print(re.findall(r'(\()|(\d+)|([-+*/])|(\))', "(1+(4+5+2)-3)+(6+8)"))
