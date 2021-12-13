# n 不会超过 9 x 10^8。


# 从 1 开始，移除所有包含数字 9 的所有整数，例如 9，19，29，……
# 这样就获得了一个新的整数数列：1，2，3，4，5，6，7，8，10，11，……
# 给定正整数 n，请你返回新数列中第 n 个数字是多少。1 是新数列中的第一个数字。

# 这些数字看起来就是 9 进制数字
# 答案就是第 n 个 9 进制数字。
class Solution:
    def newInteger(self, n: int) -> int:
        sb = []
        while n:
            div, mod = divmod(n, 9)
            n = div
            sb.append(str(mod))

        return int(''.join(sb)[::-1])


print(Solution().newInteger(9))
# 输出: 10
