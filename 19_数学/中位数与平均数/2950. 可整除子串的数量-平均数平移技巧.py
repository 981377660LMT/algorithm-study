# https://leetcode.cn/problems/number-of-divisible-substrings/description/
# 如果字符串的字符的映射值的总和可以被字符串的长度整除，则该字符串是 可整除 的。
# 给定一个字符串 s，请返回 s 的 可整除子串 的数量。
# 子串 是字符串内的一个连续的非空字符序列。
# n<=1e5


# 解:
# !可能的平均值最多为 9 个正整数 => 枚举平均数.
# 枚举每一个可能的平均值，遍历数组 nums，将数组中的每一个数减去该平均值，然后求区间和为 0 的子数组个数。

from collections import defaultdict


MAPPING = {
    "a": 1,
    "b": 1,
    "c": 2,
    "d": 2,
    "e": 2,
    "f": 3,
    "g": 3,
    "h": 3,
    "i": 4,
    "j": 4,
    "k": 4,
    "l": 5,
    "m": 5,
    "n": 5,
    "o": 6,
    "p": 6,
    "q": 6,
    "r": 7,
    "s": 7,
    "t": 7,
    "u": 8,
    "v": 8,
    "w": 8,
    "x": 9,
    "y": 9,
    "z": 9,
}


class Solution:
    def countDivisibleSubstrings(self, word: str) -> int:
        def solve(mean: int) -> int:
            preSum = defaultdict(int, {0: 1})
            res, curSum = 0, 0
            for c in word:
                curSum += MAPPING[c] - mean
                res += preSum[curSum]
                preSum[curSum] += 1
            return res

        return sum(solve(mean) for mean in range(1, 10))
