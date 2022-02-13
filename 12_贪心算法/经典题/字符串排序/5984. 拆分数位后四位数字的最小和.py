# 请你使用 num 中的 数位 ，将 num 拆成两个新的整数 new1 和 new2
# new1 和 new2 中可以有 前导 0 ，且 num 中 所有 数位都必须使用
class Solution:
    def minimumSum(self, num: int) -> int:
        digits = sorted(list(str(num)))
        num1 = digits[0] + digits[2]
        num2 = digits[1] + digits[3]
        return int(num1) + int(num2)


print(Solution().minimumSum(num=2932))
print(Solution().minimumSum(num=4009))
