# 给你一个正整数 num ，请你将它分割成两个非负整数 num1 和 num2 ，满足：

# - num1 和 num2 直接连起来，得到 num 各数位的一个排列。
#   换句话说，num1 和 num2 中所有数字出现的次数之和等于 num 中所有数字出现的次数。
# - num1 和 num2 可以包含前导 0 。
#   请你返回 num1 和 num2 可以得到的和的 最小 值。


class Solution:
    def splitNum(self, num: int) -> int:
        chars = sorted(str(num))
        a, b = chars[::2], chars[1::2]
        return int("0" + "".join(a)) + int("0" + "".join(b))
