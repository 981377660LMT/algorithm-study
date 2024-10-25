# enumerateSubstring(子串遍历)


def enumerateSubstring(arr) -> int:
    n = len(arr)
    res = 0
    for i in range(n):
        # !子串arr[i:j+1]左端点为i，这里可以做一些预处理
        ...
        for j in range(i, n):
            print(arr[i : j + 1])
            res += 1
    return res


def enumerateSubstringGroupByTargetCount(arr, target):
    """按照target的数量分组枚举子串."""
    n = len(arr)
    nexts = [0] * (n + 1)
    nexts[n] = n
    for i in range(n - 1, -1, -1):
        nexts[i] = i if arr[i] == target else nexts[i + 1]
    for i in range(n):
        j, count = i, 1 if arr[i] == target else 0
        while j != n:
            # do something using (i, j, count)
            ...
            j = nexts[j + 1]
            count += 1


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


if __name__ == "__main__":
    enumerateSubstring("abc")

    class Solution:
        # 100348. 统计 1 显著的字符串的数量-根号值域
        # https://leetcode.cn/problems/count-the-number-of-substrings-with-dominant-ones/solutions/2860181/mei-ju-by-tsreaper-x830/
        # 给你一个二进制字符串 s。请你统计并返回其中 1 显著 的子字符串的数量。
        # !如果字符串中 1 的数量 大于或等于 0 的数量的 平方，则认为该字符串是一个 1 显著 的字符串
        # !n<=4e4
        def numberOfSubstrings(self, s: str) -> int:
            from math import isqrt

            n = len(s)
            nexts = [0] * (n + 1)
            nexts[n] = n
            for i in range(n - 1, -1, -1):
                nexts[i] = i if s[i] == "0" else nexts[i + 1]

            res = 0
            upper = isqrt(n) + 1
            for i in range(n):
                j, zeros = i, 1 if s[i] == "0" else 0
                while j != n and zeros < upper:
                    # do something using (i, j, count)
                    ones = (nexts[j + 1] - i) - zeros
                    okCount = max2(0, ones - zeros * zeros + 1)
                    len_ = nexts[j + 1] - j
                    res += min2(okCount, len_)
                    j = nexts[j + 1]
                    zeros += 1
            return res
