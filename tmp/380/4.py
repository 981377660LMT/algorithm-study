from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的字符串 s 、字符串 a 、字符串 b 和一个整数 k 。

# 如果下标 i 满足以下条件，则认为它是一个 美丽下标 ：


# 0 <= i <= s.length - a.length
# s[i..(i + a.length - 1)] == a
# 存在下标 j 使得：
# 0 <= j <= s.length - b.length
# s[j..(j + b.length - 1)] == b
# |j - i| <= k
# 以数组形式按 从小到大排序 返回美丽下标。


def getNext(needle: str) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组
    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
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


def indexOfAll(longer, shorter, start=0) -> List[int]:
    """kmp O(n+m)求搜索串 `longer` 中所有匹配 `shorter` 的位置."""
    if not shorter:
        return []
    if len(longer) < len(shorter):
        return []
    res = []
    next = getNext(shorter)
    hitJ = 0
    for i in range(start, len(longer)):
        while hitJ > 0 and longer[i] != shorter[hitJ]:
            hitJ = next[hitJ - 1]
        if longer[i] == shorter[hitJ]:
            hitJ += 1
        if hitJ == len(shorter):
            res.append(i - len(shorter) + 1)
            hitJ = next[hitJ - 1]
    return res


class Solution:
    def beautifulIndices(self, s: str, a: str, b: str, k: int) -> List[int]:
        indexes1 = indexOfAll(s, a)
        indexes2 = SortedList(indexOfAll(s, b))
        res = []
        for v in indexes1:
            prev = indexes2.bisect_left(v - k)
            if prev < len(indexes2) and abs(indexes2[prev] - v) <= k:
                res.append(v)
                continue


# s = "isawsquirrelnearmysquirrelhouseohmy", a = "my", b = "squirrel", k = 15
# 输出：[16,33]
# print(indexes1, indexes2)
print(
    Solution().beautifulIndices(s="isawsquirrelnearmysquirrelhouseohmy", a="my", b="squirrel", k=15)
)
