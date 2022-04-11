from fractions import Fraction

# https://leetcode-cn.com/problems/fraction-addition-and-subtraction/solution/python-fractionlei-shi-yong-by-981377660-w5vy/


class Solution:
    def fractionAddition(self, expression: str) -> str:
        fraction = Fraction(str(eval(expression))).limit_denominator()
        return f'{fraction.numerator}/{fraction.denominator}'


print(Solution().fractionAddition("-1/2+1/2"))
print(Solution().fractionAddition("-1/2+1/2+1/3"))
