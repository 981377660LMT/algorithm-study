from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 。

# 你需要对 s 执行以下操作 任意 次：


# 选择一个下标 i ，满足 s[i] 左边和右边都 至少 有一个字符与它相同。
# 删除 s[i] 左边 离它 最近 且相同的字符。
# 删除 s[i] 右边 离它 最近 且相同的字符。
# 请你返回执行完所有操作后， s 的 最短 长度。
class Solution:
    def minimumLength(self, s: str) -> int:
        mp = defaultdict(int)
        for c in s:
            mp[c] += 1
        res = 0
        for v in mp.values():
            mod_ = v % 2
            res += 1 if mod_ == 1 else 2
        return res
