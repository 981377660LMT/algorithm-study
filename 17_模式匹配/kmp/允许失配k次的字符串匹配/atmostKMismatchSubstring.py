from typing import Any, List, Optional, Sequence, Tuple


def atmostKMismatchSubstring(longer: str, shorter: str, k: int) -> Optional[Tuple[int, int]]:
    """最多允许失配k次的子串匹配.

    返回子串的起始位置和结束位置[start, end).
    """
    if len(shorter) > len(longer):
        return None
    dp1 = zAlgo(shorter + longer)
    dp2 = zAlgo(shorter[::-1] + longer[::-1])[::-1]
    n1, n2 = len(longer), len(shorter)
    for i in range(n2, n1 + 1):
        if dp1[i] + dp2[i - 1] >= n2 - k:
            return i - n2, i
    return None


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


def zAlgo(seq: Sequence[Any]) -> List[int]:
    """z算法求字符串每个后缀与原串的最长公共前缀长度

    z[0]=0
    z[i]是s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
    """

    n = len(seq)
    z = [0] * n
    left, right = 0, 0
    for i in range(1, n):
        z[i] = max2(min2(z[i - left], right - i + 1), 0)
        while i + z[i] < n and seq[z[i]] == seq[i + z[i]]:
            left, right = i, i + z[i]
            z[i] += 1
    return z


if __name__ == "__main__":
    # 3303. 第一个几乎相等子字符串的下标(允许失配k次的最长子串)
    # https://leetcode.cn/problems/find-the-occurrence-of-first-almost-equal-substring/solutions/2934098/qian-hou-zhui-fen-jie-z-shu-zu-pythonjav-0est/
    #
    # 前后缀分解 + z数组
    class Solution:
        def minStartingIndex(self, s: str, pattern: str) -> int:
            res = atmostKMismatchSubstring(s, pattern, k=1)
            if res is None:
                return -1
            return res[0]
