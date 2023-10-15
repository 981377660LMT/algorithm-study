from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个二进制字符串 s 和一个正整数 k 。

# 如果 s 的某个子字符串中 1 的个数恰好等于 k ，则称这个子字符串是一个 美丽子字符串 。

# 令 len 等于 最短 美丽子字符串的长度。

# 返回长度等于 len 且字典序 最小 的美丽子字符串。如果 s 中不含美丽子字符串，则返回一个 空 字符串。

# 对于相同长度的两个字符串 a 和 b ，如果在 a 和 b 出现不同的第一个位置上，a 中该位置上的字符严格大于 b 中的对应字符，则认为字符串 a 字典序 大于 字符串 b 。

# 例如，"abcd" 的字典序大于 "abcc" ，因为两个字符串出现不同的第一个位置对应第四个字符，而 d 大于 c 。


class Solution:
    def shortestBeautifulSubstring(self, s: str, k: int) -> str:
        res = []
        # 枚举所有子串
        for i in range(len(s)):
            for j in range(i + 1, len(s) + 1):
                cur = s[i:j]
                if cur.count("1") == k:
                    res.append(cur)
        if not res:
            return ""
        res.sort(key=lambda x: (len(x), x))
        return res[0]
