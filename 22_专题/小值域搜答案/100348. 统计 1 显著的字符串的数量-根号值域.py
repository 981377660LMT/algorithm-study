# 100348. 统计 1 显著的字符串的数量-根号值域
# https://leetcode.cn/problems/count-the-number-of-substrings-with-dominant-ones/solutions/2860181/mei-ju-by-tsreaper-x830/
# 给你一个二进制字符串 s。
# 请你统计并返回其中 1 显著 的
# 子字符串的数量。
# !如果字符串中 1 的数量 大于或等于 0 的数量的 平方，则认为该字符串是一个 1 显著 的字符串
# !n<=4e4
#
# 小值域枚举
# !注意到子串中0的个数不是很多 => 枚举0的个数
# !枚举左端点遍历子串，再枚举子串中0的个数
# O(nsqrtn)


from math import isqrt


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def numberOfSubstrings(self, s: str) -> int:
        nums = [1 if ch == "1" else 0 for ch in s]
        n = len(nums)
        nexts = [0] * (n + 1)
        nexts[n] = n
        for i in range(n - 1, -1, -1):
            nexts[i] = i if nums[i] == 0 else nexts[i + 1]

        res = 0
        upper = isqrt(n) + 1
        for i in range(n):  # 枚举子串开头
            j, zeros = i, 1 if nums[i] == 0 else 0
            while j != n and zeros < upper:  # 从左到右枚举子串里的 0，直到数量超出限制
                ones = (nexts[j + 1] - i) - zeros  # 子串里1的数量
                okCount = max2(0, ones - zeros * zeros + 1)  # 子串右端点最多能左移几步
                len_ = nexts[j + 1] - j
                res += min2(okCount, len_)  # 和这一段的长度取 min
                j = nexts[j + 1]
                zeros += 1
        return res


print(Solution().numberOfSubstrings(s="101101"))  # 3
