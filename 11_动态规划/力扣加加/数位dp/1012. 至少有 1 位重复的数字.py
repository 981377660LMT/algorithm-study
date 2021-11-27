from math import factorial

# 所有数字-没有重复的数字
# 类似357 902
# 给定正整数 N，返回小于等于 N 且具有至少 1 位重复数字的正整数的个数。
# 1 <= N <= 10^9

# 求出小于等于4321这个数的所有重复数，中心思想就是：“因为4321是一个4位数，
# 4位数最大为9999，所以先求出小于9999这个数的不重复的数有多少，
# 之后再求出大于4321-9999这个区间的不重复数有多少，
# 最后在求这两个区间的差就可以了。”
# https://leetcode-cn.com/problems/numbers-with-repeated-digits/solution/zong-he-ge-lu-da-shen-xiang-xi-jie-ti-guo-cheng-by/


class Solution:
    def numDupDigitsAtMostN(self, N: int) -> int:
        def A(n: int, k: int) -> int:
            return 1 if k == 0 else A(n, k - 1) * (n - k + 1)
            return factorial(n) // factorial(n - k)

        s = str(N)
        n = len(s)
        no_dup = 0

        # 求出1到n-1位数的不重复数字
        for digit_len in range(1, n):
            # 第一位是1-9里选 后面是从剩下的9个数里面排列剩下的位置
            no_dup += 9 * A(9, digit_len - 1)

        # n位数时,遍历每位；
        # 判断是否是以前出现过得数字
        visited = set()
        for i in range(n):
            digit = int(s[i])
            # 依次迭代出小于当前位的有哪些
            for smaller in range(0 if i else 1, digit):
                if smaller not in visited:
                    # 从剩下的数字填充
                    no_dup += A(10 - (i + 1), n - (i + 1))
            if digit in visited:
                break
            visited.add(digit)

        # 每个数字都看过，没有break
        if n == len(set(list(s))):
            no_dup += 1

        return int(N - no_dup)


print(Solution().numDupDigitsAtMostN(9876))
print(Solution().numDupDigitsAtMostN(20))
