import re
from math import nan
from operator import add, sub, mul, floordiv

WEIGHT = dict(zip("(+-*/", [nan, 1, 1, 2, 2]))
OPT = dict(zip("+-*/", [add, sub, mul, floordiv]))


class Solution:
    def calculate2(self, s: str) -> int:
        """不能用eval python的eval有层数限制200 而js没有

        SyntaxError: too many nested parentheses
        直接eval(s)的话会报MemoryError，考虑从内到外依次计算括号
        """
        while True:
            left = s.rfind("(")
            right = s.find(")", left)
            if left == -1 or right == -1:
                break
            s = s[:left] + str(eval(s[left : right + 1])) + s[right + 1 :]
        return eval(s)

    def calculate(self, s: str) -> int:
        """基本计算器-支持+-*/和括号 (双栈+运算符优先级)"""

        def apply():
            b, a = nums.pop(), nums.pop()
            nums.append(OPT[ops.pop()](a, b))

        # 处理正负号: 开头或左括号后的 +/-
        if s.lstrip()[0] in "+-":
            s = "0" + s
        s = s.replace("(-", "(0-").replace("(+", "(0+")

        nums, ops = [], []
        for tok in re.findall(r"\d+|[+\-*/()]", s):
            if tok.isdigit():
                nums.append(int(tok))
            elif tok == "(":
                ops.append(tok)
            elif tok == ")":
                while ops[-1] != "(":
                    apply()
                ops.pop()
            else:  # +-*/
                while ops and ops[-1] != "(" and WEIGHT[ops[-1]] >= WEIGHT[tok]:
                    apply()
                ops.append(tok)
        while ops:
            apply()
        return nums[0]


if __name__ == "__main__":
    print(Solution().calculate("1 + 1"))
    print(Solution().calculate("   (  3 ) "))
    assert Solution().calculate("(1+(4+5+2)-3)+(6+8)") == 23

    print(re.split(r"(\+|\-|\*|\/)", "1+1-22"))
    print(re.findall(r"(\()|(\d+)|([-+*/])|(\))", "(1+(4+5+2)-3)+(6+8)"))
