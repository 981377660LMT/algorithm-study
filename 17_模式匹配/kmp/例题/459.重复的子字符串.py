from typing import List


def getNext(needle: str) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组

    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    https://www.ruanyifeng.com/blog/2013/05/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm.html
    """
    next = [0] * len(needle)
    j = 0

    for i in range(1, len(needle)):
        while j and needle[i] != needle[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]

        if needle[i] == needle[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1

        next[i] = j

    return next


class Solution:
    def repeatedSubstringPattern(self, s: str) -> bool:
        """给定一个非空的字符串 s ，检查是否可以通过由它的一个子串重复多次构成。"""
        return s in (s + s)[1:-1]  # !掐头去尾
        nexts = getNext(s)
        lps = nexts[-1]
        return lps != 0 and len(s) % (len(s) - lps) == 0

        n = len(s)
        for i in range(1, n // 2 + 1):
            if n % i == 0:
                if all(s[j] == s[j - i] for j in range(i, n)):
                    return True
        return False


print(Solution().repeatedSubstringPattern("abab"))
