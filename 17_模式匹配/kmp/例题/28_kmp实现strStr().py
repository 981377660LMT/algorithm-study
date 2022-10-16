from typing import List


def getNext(shorter: str) -> List[int]:
    """kmp O(n)求 `shorter`串的 `next`数组

    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    https://www.ruanyifeng.com/blog/2013/05/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm.html
    """
    next = [0] * len(shorter)
    j = 0

    for i in range(1, len(shorter)):
        while j and shorter[i] != shorter[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]
        if shorter[i] == shorter[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1
        next[i] = j

    return next


class Solution:
    def strStr(self, longer: str, shorter: str) -> int:
        nexts = getNext(shorter)
        j = 0
        for i in range(len(longer)):
            while j and longer[i] != shorter[j]:
                j = nexts[j - 1]
            if longer[i] == shorter[j]:
                j += 1
            if j == len(shorter):
                return i - len(shorter) + 1
        return -1
