from math import factorial

# 所有数字-没有重复的数字
# 类似357 902
# 给定正整数 N，返回小于等于 N 且具有至少 1 位重复数字的正整数的个数。
# 1 <= N <= 10^9

# https://leetcode.com/problems/numbers-with-repeated-digits/discuss/256725/JavaPython-Count-the-Number-Without-Repeated-Digit
class Solution:
    def numDupDigitsAtMostN(self, N: int) -> int:
        def A(n: int, k: int) -> int:
            return 1 if k == 0 else A(n, k - 1) * (n - k + 1)
            return factorial(n) // factorial(n - k)

        s = str(N)
        n = len(s)
        res = 0

        # 求出1到n-1位数的不重复数字
        # 第一位是1-9里选 后面是从剩下的9个数里面排列剩下的位置
        for digit_len in range(1, n):
            res += 9 * A(9, digit_len - 1)

        # n位数时,遍历每位；
        # 判断是否是以前出现过得数字
        # 依次迭代出小于当前位的有哪些
        visited = set()
        for i in range(n):
            digit = int(s[i])
            for smaller in range(0 if i != 0 else 1, digit):
                if smaller not in visited:
                    # 从剩下的数字填充
                    res += A(10 - (i + 1), n - (i + 1))
            if digit in visited:
                break
            visited.add(digit)

        # 每个数字都看过，没有break
        if n == len(set(list(s))):
            res += 1

        return int(N - res)


print(Solution().numDupDigitsAtMostN(9876))
print(Solution().numDupDigitsAtMostN(20))
print(Solution().numDupDigitsAtMostN(11))  # 1
