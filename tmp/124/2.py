from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 。

# 请你进行以下操作直到 s 为 空 ：


# 每次操作 依次 遍历 'a' 到 'z'，如果当前字符出现在 s 中，那么删除出现位置 最早 的该字符。
# 请你返回进行 最后 一次操作 之前 的字符串 s 。


class Solution:
    def lastNonEmptyString(self, s: str) -> str:
        mp = defaultdict(list)
        for i, v in enumerate(s):
            mp[v].append(i)
        maxLen = max(len(v) for v in mp.values())
        res = []
        for k, v in mp.items():
            if len(v) == maxLen:
                res.append((k, v[-1]))
        res.sort(key=lambda x: x[1])
        return "".join([v[0] for v in res])


# s = "aabcbbca"
print(Solution().lastNonEmptyString("aabcbbca"))
