import re

from fractions import Fraction
from math import gcd, lcm


# 592. 分数加减运算
# https://leetcode-cn.com/problems/fraction-addition-and-subtraction/solution/python-fractionlei-shi-yong-by-981377660-w5vy/


class Solution:
    def fractionAddition2(self, expression: str) -> str:
        fraction = Fraction(str(eval(expression))).limit_denominator()
        return f"{fraction.numerator}/{fraction.denominator}"

    def fractionAddition(self, expression: str) -> str:
        """常规解法"""
        if expression[0] != "-":
            expression = "+" + expression

        tokens = re.split(r"([+-])", expression)
        signs, nums = [], []
        for token in tokens:
            if token == "":
                continue
            if token in "+-":
                signs.append(1 if token == "+" else -1)
            else:
                a, b = token.split("/")
                nums.append((int(a), int(b)))

        lcm_ = lcm(*[b for _, b in nums])
        sum_ = sum(sign * a * lcm_ // b for sign, (a, b) in zip(signs, nums))
        if sum_ == 0:
            return "0/1"
        gcd_ = gcd(lcm_, sum_)
        return f"{sum_ // gcd_}/{lcm_ // gcd_}"


print(Solution().fractionAddition("-1/2+1/2"))
print(Solution().fractionAddition("-1/2+1/2+1/3"))
